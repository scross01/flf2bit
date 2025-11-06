
FIGLET_FONTS_DIR = figlet-fonts
BIT_DIR = bit

C64_FONTS = $(FIGLET_FONTS_DIR)/C64-fonts/*.flf
BDF_FONTS = $(FIGLET_FONTS_DIR)/bdffonts/*.flf

# built flf2bit
flf2bit: main.go
	go build .

$(BIT_DIR):
	git clone https://github.com/superstarryeyes/bit.git

$(FIGLET_FONTS_DIR):
	git clone https://github.com/cmatsuoka/figlet-fonts.git


# convert the Figlet C64-fonts to .bit fonts in the local fonts directory
c64fonts: flf2bit figlet-fonts
	mkdir -p fonts
	for font in $(C64_FONTS); do \
		./flf2bit --name "C64-fonts $$(basename "$${font%.flf}")" --author "Converted from figlet-fonts with flf2bit by Stephen Cross." --license "see https://github.com/cmatsuoka/figlet-fonts/blob/master/README" $$font fonts/c64-$$(basename "$${font%.flf}").bit; \
	done

# convert the Figlet bdffonts to .bit fonts in the local fonts directory
bdffonts: flf2bit figlet-fonts
	mkdir -p fonts
	for font in $(BDF_FONTS); do \
		./flf2bit --name "BDF-fonts $$(basename "$${font%.flf}")" --author "Converted from figlet-fonts with flf2bit by Stephen Cross." --license "see https://github.com/cmatsuoka/figlet-fonts/blob/master/README" $$font fonts/bdf-$$(basename "$${font%.flf}").bit; \
	done

# copy the fonts to the bit ansifonts fonts directory and rebuild bit
install: $(BIT_DIR)
	cp fonts/* $(BIT_DIR)/ansifonts/fonts
	cd $(BIT_DIR) && go build -o bit ./cmd/bit

clean:
	rm -r fonts/*.bit
	rm flf2bit

fmt:
	go fmt
	mdformat --wrap 80 README.md

.PHONY: convert install c64fonts fmt
