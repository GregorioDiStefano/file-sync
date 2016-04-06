import socket
from binascii import hexlify
import lz4
import os
import json
import hashlib
import threading
import sys
import binascii
import struct

DEBUG = False
COMPRESSED = 1 << 0
JSON_SERVER = 1 << 1
FILE_PAYLOAD = 1 << 3
EOF = 1 << 7
output_folder = "output/"


class LocalData(object):
    sending_json = ""
    existing_files = []

    def check(self, json):
        hash_sample_size = 512

        self.sending_json = json
        files = self.sending_json["FilePath"].keys()
        for f in files:
            expected_fn = self.sending_json["FilePath"][f]["FileName"]
            expected_sz = self.sending_json["FilePath"][f]["FileSize"]
            expected_fn = output_folder + expected_fn
            print "Checking hash of: %s" % (expected_fn)
            if os.path.exists(expected_fn):
                size = os.stat(expected_fn).st_size
                if size < expected_sz and size > hash_sample_size:
                    print "File: %s exists, but seems incomplete." % (expected_fn)
                    self.existing_files.append({"FileName": expected_fn,
                                                "FileSize": size,
                                                "SampleHashSize": hash_sample_size,
                                                "SampleFileHash":
                                                    MD5Check.get(
                                                        expected_fn,
                                                        last_bytes=hash_sample_size)
                                                })

                elif size == expected_sz:
                    print "File already exists and seems complete."
                    self.existing_files.append({"FileName": expected_fn,
                                                "Status": "complete"})
            else:
                    print expected_fn, " file not found locally."
            print self.existing_files

    def get_json(self):
            return json.dumps(self.existing_files)

class MD5Check(object):
    @staticmethod
    def get(filename, last_bytes=False):
        filename = filename
        hash_md5 = hashlib.md5()

        print "Checking file: ", filename
        with open(filename, "rb") as f:
            if last_bytes:
                fs = os.stat(filename).st_size
                f.seek(fs - last_bytes)
                print "Reading hash at position: ", f.tell()
            while True:
                chunk = f.read(128 * 128)
                if len(chunk):
                    hash_md5.update(chunk)
                else:
                    break
            return hash_md5.hexdigest()

    @staticmethod
    def equal(filename, expected_hash):
        assert MD5Check.get(filename) == expected_hash


class ResponseParser(object):
    ld = LocalData()

    seen_files = []
    speed_thread = None
    last_size = 0
    length_left = False
    current_key = None
    fh = None
    metadata_json = None
    buf = ""
    connection = None

    def __init__(self, connection):
        self.connection = connection

    def print_speed(self):
        self.speed_thread = threading.Timer(5.0, self.print_speed)
        self.speed_thread.start()

        sys.stdout.write(str((self.fh.tell() - self.last_size)/5.0 / (1024.0*1024.0)) +
                         "mb/sec\r")
        sys.stdout.flush()
        self.last_size = self.fh.tell()

    def new_file(self, filename):
        fixed_filename = str(output_folder + filename)

        dir_to_make = '/'.join(fixed_filename.split("/")[0:-1])
        if not os.path.exists(dir_to_make):
            os.makedirs(dir_to_make)

        if os.path.exists(fixed_filename):
            os.remove(fixed_filename)

        self.speed_thread and self.speed_thread.cancel()
        self.fh = open(fixed_filename, "a")
        self.print_speed()

    def set_metadata_json(self, data):
        print "JSON processed."

        json_data = json.loads(lz4.loads(data))
        self.metadata_json = JSONParser(json_data)
        # metadata (JSON) is read now, we are about to get the
        # first file.
        self.ld.check(self.metadata_json.json_data)

    def check_if_complete(self):
        return self.fh.tell() == self.metadata_json.key_to_size(self.current_key)

    def process_complete_file(self):
        self.fh.close()
        self.speed_thread and self.speed_thread.cancel()
        file_to_hash = output_folder + self.metadata_json.key_to_filename(self.current_key)
        expected_hash = self.metadata_json.key_to_hash(self.current_key)
        t = threading.Thread(target=MD5Check.equal, args=(file_to_hash,
                                                          expected_hash))
        t.start()

    def read_headers(self, data):
            if self.length_left > 0 or len(data) < 9:
                return False

            meta, key, length = struct.unpack("!bii", data[0:9])
            compressed = meta & COMPRESSED

            if DEBUG:
                print "data len: ", len(data), \
                      " data: ", hexlify(data[0: 32]), \
                      " meta: ", meta, \
                      " compressed: ", compressed, \
                      " file payload: ", meta & FILE_PAYLOAD, \
                      " key: ", key, \
                      " len: ", length

            if meta & FILE_PAYLOAD and \
                self.metadata_json.key_to_filename(key) not in self.seen_files:
                        fn = self.metadata_json.key_to_filename(key)
                        self.new_file(fn)
                        self.seen_files.append(fn)
                        print "new file: ", fn

            self.current_key = key
            self.current_meta = meta

            self.length_left = length
            self.process_data(compressed, data[9:])

    def process_data(self, compressed, data):
        if len(data) >= self.length_left:
                tmp = self.length_left
                self.length_left = False
                self.buf = data[tmp:]

                if self.current_meta & JSON_SERVER:
                    self.set_metadata_json(data[0:tmp])
                    self.connection.send_data(self.ld.get_json())
                    self.buf = self.buf[9:]
                else:
                    if compressed:
                        self.fh.write(bytearray(lz4.loads(data[0:tmp])))
                    else:
                        self.fh.write(bytearray(data[0:tmp]))

                    if self.check_if_complete():
                        self.process_complete_file()

                #TODO: this is wrong
                #if len(self.buf) >= 9:
                    self.read_headers(self.buf)
        else:
            self.length_left = False

    def check(self):
        self.read_headers(self.buf)

class JSONParser(object):
    json_data = ""

    def __init__(self, json):
            self.json_data = json

    def files_to_be_sent(self):
        files = []
        for k in self.json_data["FilePath"].keys():
            files.append(self.json_data["FilePath"][k]["FileName"])
        return files

    def key_to_filename(self, key):
        key = str(key)
        return self.json_data["FilePath"][key]["FileName"]

    def key_to_hash(self, key):
        key = str(key)
        return self.json_data["FilePath"][key]["FileHash"]

    def key_to_size(self, key):
        key = str(key)
        return self.json_data["FilePath"][key]["FileSize"]

class Connection(object):
    rp = None
    s = socket.socket()
    s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    connection = None
    port = 8081
    s.bind(("0.0.0.0", port))
    s.listen(1)

    def wait_for_connection(self):
        while True:
            self.connection, addr = self.s.accept()
            if self.connection:
                break

        self.rp = ResponseParser(self)
        print "Connected."
        return True

    def read_data(self):
        while True:
            data = self.connection.recv(10 * 1024 * 1024)

            if data:
                self.rp.buf += data
                self.rp.check()
            elif not data:
                break

        if len(self.rp.buf):
            self.rp.read_headers(self.rp.buf)

        if self.rp.speed_thread:
            self.rp.speed_thread.cancel()

        assert sorted(self.rp.seen_files) == sorted(self.rp.metadata_json.files_to_be_sent())
        sys.exit(0)

    def send_data(self, payload):
        def prepare_payload(self):
            prefix = "%02x%08x%08x" % (FILE_PAYLOAD, 0x00, len(payload))
            return binascii.a2b_hex(prefix) + payload

        if payload:
            self.connection.sendall(prepare_payload(payload))
            print "Sent: ", prepare_payload(payload)

c = Connection()
c.wait_for_connection()
c.read_data()
