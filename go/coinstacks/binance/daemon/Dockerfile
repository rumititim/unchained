FROM ubuntu:jammy

RUN apt-get update -y \
  	&& apt-get install -y curl jq unzip wget \
  	&& apt-get clean \
  	&& rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ARG VERSION

RUN wget https://github.com/bnb-chain/node/releases/download/${VERSION}/linux_binary.zip && \
	unzip linux_binary.zip && \
	chmod +x linux/* && \
	mv linux/* /usr/local/bin && \
	rm -rf linux_binary.zip linux

RUN wget https://github.com/bnb-chain/node/releases/download/${VERSION}/mainnet_config.zip && \
	unzip mainnet_config.zip && \
	mv asset/mainnet/* / && \
	rm mainnet_config.zip

EXPOSE 26656 26657 26658

ENTRYPOINT ["bnbchaind"]