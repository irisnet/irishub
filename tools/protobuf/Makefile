build:
	g++ -o Json2Pb.o -c Json2Pb.cc -lprotobuf -ljsoncpp -lstdc++ -std=c++11
	g++ -o irisnet.o -c irisnet.cc -lprotobuf -ljsoncpp -lstdc++ -std=c++11
	ar r libiservice.so Json2Pb.o irisnet.o

clean:
	-rm -f *.o *.a *.so
