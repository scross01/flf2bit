
FIGLET_FONTS_DIR = figlet-fonts
BIT_DIR = bit

C64_FONTS = $(FIGLET_FONTS_DIR)/C64-fonts/*.flf
BDF_FONTS = $(FIGLET_FONTS_DIR)/bdffonts/*.flf

JAVE_FONTS = $(FIGLET_FONTS_DIR)/jave/*.flf
FIGLET_FONTS = $(FIGLET_FONTS_DIR)/ours/*.flf
TOILET_FONTS = $(FIGLET_FONTS_DIR)/toilet/*.tlf
CONTRIBUTED_FONTS = $(FIGLET_FONTS_DIR)/contributed/*.flf
MSDOS_FONTS = $(FIGLET_FONTS_DIR)/ms-dos/*.flf


# built flf2bit
flf2bit: main.go
	go build .

$(BIT_DIR):
	git clone https://github.com/superstarryeyes/bit.git

$(FIGLET_FONTS_DIR):
	git clone https://github.com/cmatsuoka/figlet-fonts.git

all-fonts: bdffonts c64fonts contributed jave msdos ours toilet

# convert the Figlet figlet-fonts/bdffonts
bdffonts: flf2bit figlet-fonts
	mkdir -p fonts
	for font in $(BDF_FONTS); do \
		./flf2bit --map-chars "#█" --name "Figlet-BDF $$(basename "$${font%.flf}")" --author "Converted from figlet-fonts with flf2bit by Stephen Cross." --license "see https://github.com/cmatsuoka/figlet-fonts/blob/master/README" $$font fonts/bdf-$$(basename "$${font%.flf}").bit; \
	done

# convert the figlet-fonts/C64-fonts
c64fonts: flf2bit figlet-fonts
	mkdir -p fonts
	for font in $(C64_FONTS); do \
		./flf2bit --map-chars "#█" --name "Figlet-C64 $$(basename "$${font%.flf}")" --author "Converted from figlet-fonts with flf2bit by Stephen Cross." --license "see https://github.com/cmatsuoka/figlet-fonts/blob/master/README" $$font fonts/c64-$$(basename "$${font%.flf}").bit; \
	done

# convert the figlet-fonts/contributed fonts
contributed: flf2bit figlet-fonts
	mkdir -p fonts
	for font in $(CONTRIBUTED_FONTS); do \
		./flf2bit --name "Figlet-Contrib $$(basename "$${font%.flf}")" --author "Converted from figlet-fonts with flf2bit by Stephen Cross." --license "see https://github.com/cmatsuoka/figlet-fonts/blob/master/README" $$font fonts/contributed-$$(basename "$${font%.flf}").bit; \
	done

# convert the figlet-fonts/jave fonts
jave: flf2bit figlet-fonts
	mkdir -p fonts
	for font in $(JAVE_FONTS); do \
		./flf2bit --name "Figlet-JavE $$(basename "$${font%.flf}")" --author "Converted from figlet-fonts with flf2bit by Stephen Cross." --license "see https://github.com/cmatsuoka/figlet-fonts/blob/master/README" $$font fonts/jave-$$(basename "$${font%.flf}").bit; \
	done

# convert the figlet-fonts/ms-dos fonts
msdos: flf2bit figlet-fonts
	mkdir -p fonts
	for font in $(MSDOS_FONTS); do \
		./flf2bit --name "Figlet-MSDOS $$(basename "$${font%.flf}")" --author "Converted from figlet-fonts with flf2bit by Stephen Cross." --license "see https://github.com/cmatsuoka/figlet-fonts/blob/master/README" $$font fonts/msdos-$$(basename "$${font%.flf}").bit; \
	done

# convert the figlet-fonts/ours fonts
ours: flf2bit figlet-fonts
	mkdir -p fonts
	for font in $(FIGLET_FONTS); do \
		./flf2bit --name "Figlet-Std $$(basename "$${font%.flf}")" --author "Converted from figlet-fonts with flf2bit by Stephen Cross." --license "see https://github.com/cmatsuoka/figlet-fonts/blob/master/README" $$font fonts/figlet-$$(basename "$${font%.flf}").bit; \
	done

# convert the figlet-fonts/toilet fonts
toilet: flf2bit figlet-fonts
	mkdir -p fonts
	for font in $(TOILET_FONTS); do \
		./flf2bit --name "Toilet $$(basename "$${font%.tlf}")" --author "Converted from figlet-fonts with flf2bit by Stephen Cross." --license "see https://github.com/cmatsuoka/figlet-fonts/blob/master/README" $$font fonts/toilet-$$(basename "$${font%.tlf}").bit; \
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

.PHONY: convert install c64fonts fmt clean all-fonts
