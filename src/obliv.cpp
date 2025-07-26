#include <iostream>

#include "server.cpp"
using namespace std;

void
intro(void)
{
	cout << "obliv - web server panel" << endl;
	cout << "https://github.com/radh1tya/obliv" << endl;
}
int
main(void)
{
	intro();
	start_socket();
	create_address();
}
