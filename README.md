# genFact
Customizable PDF invoice generator
## Installation
Clone the repository
```
git clone https://github.com/MISELLO/genFact/
```
Not mandatory but necessary if you want to update the EXIF information (the app tells you the exact command you should use)
```
sudo apt install libimage-exiftool-perl
```
## Usage
To create a default invoice:
```
git clone https://github.com/MISELLO/genFact/
cd genFact
go run genFact.go
```
Or use a file input:
```
go run genFact.go default.tsv
```
Feel free to copy and modify this .tsv file as per your needs. Also, you can do the same with the static_data/static_default.dat in order to put the information you need.
Both are text files very easy to understand.
