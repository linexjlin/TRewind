# TRewind

Text based memory rewind tool.

# Usage

1. Unzip the software to any folder and run TRewind. Then you will see the app running in the system tray.

![](https://ipfs.ee/ipfs/QmSCVxRgcFhmE5JVkdL2eF88Uc84kxz3MrNYe6u4ux2g7M/gIgcr.png)

2. Copy the content you want to remember, then click the app icon in the system tray to add the content to the system. There are 3 categories of memory: docs, ideas, web excerpts. 
    - *Add All* means index all content of the clipboard.
    - *Add First Line* means only the first line of the content will be indexed.

![](https://ipfs.ee/ipfs/QmdHZE9a7ZdcjdBg5LQorjdZfz38ytAEqoRcgVhXdaNtDF/2AH8a.png)

3. Recall the memory. Click the search menu and the app will open the browser to display a simple web UI for searching or managing the memory.

![](https://ipfs.ee/ipfs/QmTn1JeZeTPrgywS2Qejkp1WYpZ3BJwpefREBSSeCRCfgh/FSSc4.png)


# About Configuration

All configurations are in the env.txt file.

```
OLLAMA_EMBED_URL = "https://oapi.chinatcc.com/api"  # By default, this app uses a hosted Ollama server. You can change to your local Ollama server by modifying this line to `OLLAMA_EMBED_URL = "http://127.0.0.1:13134/api"`
OLLAMA_EMBED_MODE = "qwq5km"  # By default, I am using gte-qwen2-7b as my embedding mode
EMBED_DIR = "docs_db"  # Local embedding store
DEFAULT_COLLECTION = "docs"  # The default category of memory
COLLECTIONS = "idea;code;web_excerpt"  # Add/delete categories of memory
API_LISTEN_ADDR = "127.0.0.1:8601"  # Local API server listen address
REMOTE_API_ADDR="127.0.0.1:8601" # Use a remote server instead of a local API server to allow cross-platform add and recall.
```

# Credit

- [github.com/philippgille/chromem-go](https://github.com/philippgille/chromem-go)

- [golang.design/x/clipboard](https://github.com/golang-design/clipboard)

- [github.com/getlantern/systray](https://github.com/getlantern/systray)
