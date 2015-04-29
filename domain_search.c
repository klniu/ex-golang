#include <stdio.h>
#include <sys/ioctl.h>
#include <unistd.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <fcntl.h>
#include <memory.h>
#include <signal.h>
#include <netdb.h>
#include <time.h>
#include <linux/soundcard.h>

int open_socket(char *host, int clientPort);


int main(int argc,char * argv[])
{
char buffer[50000];
int fd,ret,loop;

	if (argc <= 2)
		{
		puts("Usage: domain_search <domain name> <TLD> <TLD> ... ");
		exit (0);
		}
	fd = open_socket("www.uwhois.com",999);
	
	for(loop=1;loop<argc;loop++)
		{
		strcat(buffer,argv[loop]);
		strcat(buffer," ");
		}
	strcat(buffer,"\n");

	ret = write(fd,buffer,strlen(buffer));

	while((ret = read(fd,buffer,sizeof(buffer)-1)) > 0) 
		{
		if (write(1,buffer,ret) != ret) perror("write");
		fflush(stdout);
		}
	shutdown(fd,1);
	close(fd);

	return 0;
}



int open_socket(char *host, int clientPort)
{
    int sock;
    unsigned long inaddr;
    struct sockaddr_in ad;
    struct hostent *hp;

    memset(&ad, 0, sizeof(ad));
    ad.sin_family = AF_INET;

    inaddr = inet_addr(host);
    if (inaddr != INADDR_NONE)
        memcpy(&ad.sin_addr, &inaddr, sizeof(inaddr));
    else
		{
        if ((hp = gethostbyname(host))==NULL) return -1;
        memcpy(&ad.sin_addr, hp->h_addr, hp->h_length);
		}
    ad.sin_port = htons(clientPort);
    
    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) return sock;
    
    if (connect(sock, (struct sockaddr *) &ad, sizeof(ad)) < 0)
		{
		shutdown(sock,2);
		close(sock);
		return -1;
		}

    return sock;
}
