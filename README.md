# configwalker

Configwalker is a configuration extractor for Netwalker Ransomware. It is capable of decrypting the RC4 encrypted Resource File and extracting the Ransomnote template. By default it will dump the results to disk, but you can also choose to print the config to stdout only by appending ```--print``` to the command.


## Usage

```go run path/to/sample.exe [--print]```


## Screenshots

![Configwalker Screenshot](screenshots/sc.png)


## Background Info

Not all samples of Netwalker Ransomware will contain a full configuration file. Below you can see two screenshots of Resource Hacker, one with an long encrypted config and one with a short plaintext config. This tool only works for the former format.

![long encrypted](screenshots/encrypted-config.png)


![short plaintext](screenshots/plain-config.png)


## Testing

The tool has been confirmed to successfully extract the configuration from the following samples identifiable by their SHA-256 hashsums. If you encounter an error with configwalker please file a bug report as an issue. On some occasions the Netwalker config files contain malformed json objects (fix WIP).

| SHA-256 Hashsums|
|:---------------:| 
| 1707f8647515bc7a686e7aed380ab06dd6b853b908ae98252c1e8eefa1e1d540 |
| 5d869c0e077596bf0834f08dce062af1477bf09c8f6aa0a45d6a080478e45512 |
| 46dbb7709411b1429233e0d8d33a02cccd54005a2b4015dcfa8a890252177df9 |
| 4f7bdda79e389d6660fca8e2a90a175307a7f615fa7673b10ee820d9300b5c60 |
| 27319e75c23693399977e92b9a7ba5680a7a9db448f93b3221840c61301604d5 |
