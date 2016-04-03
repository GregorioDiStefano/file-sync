import socket
from binascii import hexlify
import lz4
import os
import json
import hashlib
import threading
import sys

DEBUG = True
JSON = 0x00
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
            if os.path.exists(expected_fn):
                size = os.stat(expected_fn).st_size
                if size < expected_sz:
                    print "File: %s exists, but seems incomplete." % (expected_fn)
                    self.existing_files.append({"FileName": expected_fn,
                                                "FileSize": size,
                                                "SampleHashSize": hash_sample_size,
                                                "SampleFileHash":
                                                    MD5Check.get(
                                                        expected_fn,
                                                        last_bytes=hash_sample_size)
                                                })
                else:
                    print "File already exists and seems complete."
            else:
                    print expected_fn, " file not found locally."
            print self.existing_files

    def get_json(self):
        if self.existing_files:
            return json.dumps(self.existing_files)
        else:
            print "No existing files info set.."

class MD5Check(object):
    @staticmethod
    def get(filename, last_bytes=False):
        filename = filename
        hash_md5 = hashlib.md5()

        with open(filename, "rb") as f:
            if last_bytes:
                fs = os.stat(filename).st_size
                f.seek(fs - last_bytes)
                print "Reading hash at position: ", f.tell()
            for chunk in iter(lambda: f.read(4096), b""):
                hash_md5.update(chunk)
            return hash_md5.hexdigest()

    @staticmethod
    def equal(filename, expected_hash):
        assert MD5Check.get(filename) == expected_hash


class ResponseParser(object):
    ld = LocalData()

    speed_thread = None
    last_size = 0
    length_left = False
    current_key = None
    fh = None
    metadata_json = ""
    buf = ""
    connection = None

    def __init__(self, connection):
        self.connection = connection

    def print_speed(self):
        self.speed_thread = threading.Timer(1.0, self.print_speed)
        self.speed_thread.start()

        sys.stdout.write(str((self.fh.tell() - self.last_size) / (1024*1024)) +
                         "mb/sec\r")
        sys.stdout.flush()
        self.last_size = self.fh.tell()

    def new_file(self, filename):
        if self.fh and not self.fh.closed:
            self.fh.close()

        fixed_filename = str(output_folder + filename)

        if os.path.exists(fixed_filename):
            os.remove(fixed_filename)

        self.speed_thread and self.speed_thread.stop()
        self.fh = open(fixed_filename, "a")
        self.print_speed()

    def set_metadata_json(self):
        self.metadata_json = json.loads(lz4.loads(self.metadata_json))

    def key_to_filename(self, key):
        key = str(key)
        return self.metadata_json["FilePath"][key]["FileName"]

    def key_to_hash(self, key):
        key = str(key)
        return self.metadata_json["FilePath"][key]["FileHash"]

    def read_headers(self, data):
            if self.length_left > 0:
                return False

            compressed = int(hexlify(data[0:1]), 16)
            key = int(hexlify(data[1:5]), 16)
            length = int(hexlify(data[5:9]), 16)

            # metadata (JSON) is read now, we are about to get the
            # first file.
            if key == 1 and self.current_key == JSON:
                self.set_metadata_json()
                self.ld.check(self.metadata_json)
                existing_files_info = self.ld.get_json()
                self.connection.send_data(existing_files_info)
                sys.exit(1)

            if key != self.current_key and key > JSON:
                if key == self.current_key:

                    self.speed_thread and self.speed_thread.stop()
                    MD5Check.equal(output_folder + self.key_to_filename(key),
                                   self.key_to_hash(key))

                self.new_file(self.key_to_filename(key))
                print "new file: ", self.key_to_filename(key)

            self.current_key = key

            if DEBUG:
                print "data len: ", len(data), \
                      " data: ", hexlify(data[0: 10]), \
                      " compressed: ", compressed, \
                      " key: ", key, \
                      " len: ", length

            self.length_left = length
            self.process_data(compressed, data[9:])

    def process_data(self, compressed, data):
        if len(data) >= self.length_left:
                tmp = self.length_left
                self.length_left = False
                self.buf = data[tmp:]

                if self.current_key == JSON:
                    self.metadata_json += data[0:tmp]

                elif compressed:
                    self.fh.write(bytearray(lz4.loads(data[0:tmp])))
                else:
                    self.fh.write(bytearray(data[0:tmp]))
        else:
            self.length_left = False

    def check(self):
        self.read_headers(self.buf)


class Connection(object):
    rp = None
    s = socket.socket()
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
            data = self.connection.recv(1024 * 1024 * 10)

            if data:
                self.rp.buf += data
                self.rp.check()
            elif not data:
                break

        if len(self.rp.buf):
            self.rp.read_headers(self.rp.buf)

        self.rp.speed_thread.cancel()
        sys.exit(0)

    def send_data(self, payload):
        self.connection.send(payload)
        print "Sent: ", payload

c = Connection()
c.wait_for_connection()
c.read_data()
