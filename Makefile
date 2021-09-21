VERSION?=dev
GOOS?=linux
BUILD_PATH?=./target

build:
	rm -f $(BUILD_PATH)/esher-notifier
	go build $(GO_TAGS) $(GO_LDFLAGS) -o $(BUILD_PATH)/esher-notifier ./
	cp -v config.yml $(BUILD_PATH)/
	cp -v client_secret.json $(BUILD_PATH)/
	mkdir -p $(BUILD_PATH)/mailer/
	cp -rv mailer/mails $(BUILD_PATH)/mailer/
	cp run.sh $(BUILD_PATH)/

clean:
	rm -rf $(BUILD_PATH)
