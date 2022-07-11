cd ~/webOS_apps
ares-package first-app/

ares-install --device tv0 --remove com.aloalo.app.firstapp
ares-install --device tv1 --remove com.aloalo.app.firstapp
ares-install --device tv2 --remove com.aloalo.app.firstapp
ares-install --device tv3 --remove com.aloalo.app.firstapp
ares-install --device tv4 --remove com.aloalo.app.firstapp
ares-install --device tv5 --remove com.aloalo.app.firstapp

ares-install --device tv0 com.aloalo.app.firstapp_0.0.1_all.ipk
ares-install --device tv1 com.aloalo.app.firstapp_0.0.1_all.ipk
ares-install --device tv2 com.aloalo.app.firstapp_0.0.1_all.ipk
ares-install --device tv3 com.aloalo.app.firstapp_0.0.1_all.ipk
ares-install --device tv4 com.aloalo.app.firstapp_0.0.1_all.ipk
ares-install --device tv5 com.aloalo.app.firstapp_0.0.1_all.ipk


ares-launch --device tv0 com.aloalo.app.firstapp
ares-launch --device tv1 com.aloalo.app.firstapp
ares-launch --device tv2 com.aloalo.app.firstapp
ares-launch --device tv3 com.aloalo.app.firstapp
ares-launch --device tv4 com.aloalo.app.firstapp
ares-launch --device tv5 com.aloalo.app.firstapp
