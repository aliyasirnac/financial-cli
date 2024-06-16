#!/bin/sh

# CLI uygulamasının bulunduğu dizini belirleyin
PROJECT_DIR=$(pwd)

# CLI uygulamasının yürütülebilir dosyasını /usr/local/bin dizinine kopyalayın
sudo cp "$PROJECT_DIR/bin/nachboard" /usr/local/bin/nachboard

# ZSHRC dosyasına alias ekleyin
echo 'alias nachboard="/usr/local/bin/nachboard"' >> ~/.zshrc

# ZSHRC dosyasını yeniden yükleyin
source ~/.zshrc

echo "CLI application has been set up. You can now run it using the 'nachboard' command."