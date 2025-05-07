stage_2_persistent_path="/bin/xiftooling"
stage_2_service_name="crash-reporter"

echo "[+] Checking for crontab..."
crontab -l 2>/dev/null 1>/dev/null
if [ $? -eq 0 ]; then
    echo "[+] Crontab found. Removing entry..."
    (crontab -l | grep -v "$stage_2_persistent_path") | crontab - 
fi

echo "[+] Checking for systemd..."
systemctl --version 2>/dev/null 1>/dev/null
if [ $? -eq 0 ]; then
    service_path="/etc/systemd/system/$stage_2_service_name.service"
    echo "[+] Systemd found. removing service at $service_path..."
    systemctl disable --now "$stage_2_service_name.service"
    rm -f $service_path
    systemctl daemon-reload
fi

echo "[+] Checking for upstart..."
initctl version 2>/dev/null 1>/dev/null
if [ $? -eq 0 ]; then
    service_path="/etc/init/$stage_2_service_name.conf"
    echo "[+] Upstart found. Removing service at $service_path..."
    initctl stop "$stage_2_service_name"
    rm -f $service_path
    initctl reload-configuration
fi

echo "[+] Checking for init.d..."
service -h 2>/dev/null 1>/dev/null
if [ $? -eq 0 ]; then
    service_path="/etc/init.d/$stage_2_service_name"
    echo "[+] Init.d found. Removing service at $service_path..."
    service $stage_2_service_name stop
    update-rc.d -f $stage_2_service_name remove
    rm -f $service_path
fi

echo "[+] Checking for rc.local..."
if [ -f /etc/rc.local ]; then
    echo "[+] rc.local found. Removing entry..."
    sed -i /"$stage_2_persistent_path"/d /etc/rc.local
fi

echo "[+] Checking for bashrc..."
if [ -f ~/.bashrc ]; then
    echo "[+] bashrc found. Removing entry..."
    sed -i /"$stage_2_persistent_path"/d ~/.bashrc # Try to remove the entry
fi

echo "[+] Removing $stage_2_persistent_path..."
rm -f $stage_2_persistent_path

echo "[+] Done!"