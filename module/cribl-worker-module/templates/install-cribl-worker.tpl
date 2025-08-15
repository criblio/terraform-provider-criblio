#!/bin/bash
set -e

# Cribl Edge Installation Script
# Based on official documentation: https://docs.cribl.io/edge/deploy-single-instance/

echo "Starting Cribl Edge installation..."

# Check minimum requirements (based on Cribl docs)
echo "Checking system requirements..."

# Verify we have curl
if ! command -v curl &> /dev/null; then
    echo "Installing curl..."
    if command -v apt-get &> /dev/null; then
        sudo apt-get update && sudo apt-get install -y curl
    elif command -v yum &> /dev/null; then
        sudo yum install -y curl
    elif command -v dnf &> /dev/null; then
        sudo dnf install -y curl
    fi
fi

# Install Cribl Edge using official installation URL
echo "Installing Cribl Node..."
curl -fsSL '${cribl_install_url}' | bash -

# Verify installation
if [ -d "/opt/cribl" ]; then
    echo "Cribl Node installed successfully in /opt/cribl"
    
    # Set proper permissions (following Cribl best practices)
    sudo chown -R cribl:cribl /opt/cribl 2>/dev/null || true
    
    # Enable auto-start on boot (optional - customers can customize)
    # sudo /opt/cribl/bin/cribl boot-start enable
    
    echo "Cribl Node installation completed successfully!"
else
    echo "Warning: Cribl installation directory not found. Please verify installation."
fi

# Add your custom installation steps below this line
# Examples based on Cribl documentation:
# - Configure Cribl Edge settings
# - Set up SSL/TLS certificates  
# - Configure data sources and destinations
# - Set environment variables
# - Install additional monitoring tools
# - Configure firewall rules for required ports
# - Set up log rotation
# - Configure backup procedures 