/**
  socket -> bind -> listen -> accept -> read/write
  **/

#include <iostream>
#include <sys/socket.h>
#include <arpa/inet.h>

using namespace std;

#define PORT 2305

int
start_socket()
{
	int sockfd = socket(AF_INET, SOCK_STREAM, 0);
	if (sockfd == -1) {
		cout << "error creating socket" << endl;
		exit(1);
	}
	return sockfd;
}

void
create_address()
{
	int sockfd = start_socket();

	struct sockaddr_in host_addr;
	int host_addrlen = sizeof(host_addr);

	host_addr.sin_family = AF_INET;
	host_addr.sin_port = htons(PORT);
	host_addr.sin_addr.s_addr = htonl(INADDR_ANY);

	// bind socket :)
	if (bind(sockfd, (struct sockaddr *)&host_addr, host_addrlen) != 0) {
		cout << "error webserver (bind)" << endl;
		exit(1);
	}

	cout << "socket successfully bound to address" << endl;
}
