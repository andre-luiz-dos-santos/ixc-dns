all: build_bin build_root build_deb
	rm -rf root

build_bin:
	cd .. && go build

build_root:
	rm -rf root
	mkdir -p root/opt/ixc-dns/bin
	cp ../ixc-dns root/opt/ixc-dns/bin/
	mkdir -p root/opt/ixc-dns/etc
	cp ../ixc-dns.toml root/opt/ixc-dns/etc/ixc-dns.toml

build_deb:
	~/.gem/ruby/*/bin/fpm \
		--input-type dir \
		--output-type deb \
		--name ixc-dns \
		--version 1.0 \
		--package ixc-dns.deb --force \
		--config-files opt/ixc-dns/etc \
		--directories opt/ixc-dns \
		--after-install after-install.sh \
		--deb-systemd ixc-dns.service \
		--chdir root .

install_fpm:
	gem install --no-document --user-install fpm
