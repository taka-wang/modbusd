# v0.2.2

## Done
- [x] rewrite test server in c, test cases in golang
- [x] update travis ci with new test env
- [x] implement timeout.set, timeout.get functions
- [x] implement json_get_double, json_set_double, json_get_long
- [x] refactor tid data type


---

# v0.2.1

## Done
- [x] fix uthash sizeof issue (1000 items); cause: memset to the wrong size 
- [x] fix char pointer key issue (hash); use char array instead of char pointer (unkown length)
- [x] modbus_connect hang issue; set tcp timeout
- [x] handle 'reset by peer' issue; workaround: set connection flag to false :warning:
- [x] implement keep connection mechanism via hashtable
- [x] implement FC (1~6, 15, 16)
- [x] assign daemon version number from the latest git tag
- [x] implement syslog and flag mechanism
- [x] implement read/write config mechanism
- [x] define default config
- [x] implement dummy modbus server in node.js for testing
- [x] support ipv4/v6 ip address string
- [x] refactor int port to char * port
- [x] support docker compose
- [x] support valgrind (disable now)
- [x] implement set/get timeout command
- [x] support mocha and async test (runs slow on cloud server)
- [x] add versioneye support (depends check)
- [x] support armhf
- [x] deploy doxygen document automatically

---

## TODO List

- [ ] enhance reconnect mechanism :clap:
- [ ] refine field name for psmb

