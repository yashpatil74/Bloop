# Bloop - P2P File Transfer System

Bloop is a **peer-to-peer (P2P) file sharing system** designed for seamless file transfers between devices on the same network. It uses **WebSockets** for real-time communication and **TCP sockets** for efficient file transfer, leveraging **mDNS** for automatic device discovery.

## Features 🚀

- **Real-time Device Discovery** - Uses mDNS to detect nearby nodes.
- **WebSocket Communication** - Establishes real-time messaging between nodes.
- **File Transfer via TCP** - Secure and fast file transfers.
- **Beautiful UI** - Modern and user-friendly interface.
- **Automatic File Storage** - Received files are saved in `~/Downloads/bloop/`.

---

## ⚠️ Important: Windows-Only Support (For Now)

Currently, Bloop **only works properly on Windows** due to **mDNS discovery issues on MacOS and Linux**.

I am working on cross-platform support, but for now, it functions best on **Windows 10/11**.

**⚠️ WARNING:** This project is in early development and **should not be used in public networks**. It lacks full security implementations such as encryption.

---

## Installation 🔧

### 1. Clone the Repository

```sh
git clone https://github.com/yashpatil74/bloop.git
cd bloop
```

### 2. Install Dependencies

#### **Backend (Go)**

Make sure you have **Go 1.18+** installed.

```sh
go mod tidy
```

#### **Frontend (Next.js)**

Ensure **Node.js 16+** and **npm** or **yarn** are installed.

```sh
cd web/next
npm install  # or yarn install
```

### 3. Build and Run the Application

#### **Build the Frontend**

```sh
cd web/next
npm run build
```

#### **Build and Run the Final App (Frontend + Backend)**

```sh
cd cmd
go run main.go
// or
go build -o bloop.exe
```

Once started, open **[http://localhost:8787/app](http://localhost:8787/app)** in your browser.

---

## Usage 🖥️

### **1. Detect Nearby Devices**

- When you start the app, it will **automatically scan** for other Bloop nodes on the network.
- A list of available devices will be displayed.

### **2. Sending a File**

1. Select a device from the list.
2. Choose a file to send.
3. Click **Send File**.
4. The recipient will receive a request to accept or decline the file.

### **3. Receiving a File**

- A pop-up will appear when another device sends you a file.
- Click **Accept** to start the file transfer.
- The file will be saved in `~/Downloads/bloop/`.

---

## Troubleshooting 🛠️

### **Firewall Issues**

Make sure **port 8787 (WebSocket) and 9090 (TCP)** are open on your firewall.

### **WebSocket Connection Issues**

If the frontend isn't connecting to the backend:

- Check if the **Go backend is running**.
- Verify that **WebSockets are allowed** in the browser.

---

## Future Improvements 🏗️

- **End-to-End Encryption (E2EE) for Secure Transfers** 🔒
- **Progress Bar & Transfer Speeds** 📊
- **Background File Transfers** 🎯
- **Cross-platform Support (MacOS & Linux)** 🏗️
- **Mobile App (iOS & Android)** 📱

---

## Contributing 🤝

Want to improve Bloop? Follow these steps:

1. Fork the repository.
2. Create a new branch (`feature/amazing-feature`).
3. Submit a pull request!

---

## Acknowledgments 💡

- [Go WebSockets](https://github.com/gorilla/websocket) 🦍
- [mDNS Service](https://github.com/grandcat/zeroconf) 🌎
- [Next.js for UI](https://nextjs.org/) ⚛️

---

### **Made with ❤️ by Yash**

