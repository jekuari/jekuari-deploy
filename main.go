/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

/*
	script to run the container to compile to alpine arm64

	docker run --name jkdpb -d -it $(docker build -q -t generate-nginx-block-compiler .) && docker container cp jkdpb:/bin/generate-nginx-block ./generate-nginx-block && docker container rm -f jkdpb && scp generate-nginx-block AWS-EC2:~/bin/

	script to copy the binary from the container to the host

	docker container cp jkdpb:/bin/generate-nginx-block ./generate-nginx-block

	script to send the binary to the server

	scp generate-nginx-block AWS-EC2:~/bin/

*/

func main() {

	fmt.Println("nginx server block generator")

	var domain string

	if (len(os.Args) > 1) && (os.Args[1] == "-d") {

		if len(os.Args) > 2 {
			domain = os.Args[2]
		}
	} else {

		fmt.Println("Enter the domain name")

		fmt.Scanln(&domain)
	}

	interpolatedString := fmt.Sprintf(
		`
		server {
			listen 80;
			listen [::]:80;

			root /var/www/%s/html;
			index index.html index.htm index.nginx-debian.html;

			server_name %s www.%s;

			location / {
					  try_files $uri $uri/ /index.html;
			}
		}`, domain, domain, domain)

	os.MkdirAll("/etc/nginx/sites-available", 0755)

	file, err := os.Create("/etc/nginx/sites-available/" + domain)

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	_, err2 := file.WriteString(interpolatedString)

	if err2 != nil {
		fmt.Println(err2)
	}

	fmt.Println("File created: ", domain)

	// create a symlink

	os.MkdirAll("/etc/nginx/sites-enabled", 0755)

	err3 := os.Symlink("/etc/nginx/sites-available/"+domain, "/etc/nginx/sites-enabled/"+domain)

	if err3 != nil {
		fmt.Println(err3)
	}

	fmt.Println("Symlink created: ", domain)

	certbotCommand := fmt.Sprintf("certbot --nginx -d %s -d www.%s", domain, domain)

	errOnRunningCertbot := exec.Command("sh", certbotCommand).Run()

	if errOnRunningCertbot != nil {
		fmt.Println(errOnRunningCertbot, "Certbot failed")
	}

	fmt.Println("Certbot ran successfully")

	// restart nginx

	errOnRestartingNginx := exec.Command("sh", "systemctl restart nginx").Run()

	if errOnRestartingNginx != nil {
		fmt.Println(errOnRestartingNginx, "Nginx failed to restart")
	}

	fmt.Println("Nginx restarted successfully")

	for {
		fmt.Println("Do you want to add another domain? (y/n)")

		var answer string

		fmt.Scanln(&answer)

		if answer == "y" {
			main()
		} else if answer == "n" {
			os.Exit(0)
		} else {
			fmt.Println("Invalid input")
		}
	}
}
