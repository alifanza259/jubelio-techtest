
var conn;

async function init() {
    if (sessionStorage.getItem('userId') == null) {
        window.location.href = 'homepage.html'
    }

    document.getElementById("my-id").innerHTML = "Your ID: " + sessionStorage.getItem("userId")

    window.onload = function () {
        const accessToken = sessionStorage.getItem("accessToken")
        conn = new WebSocket(`ws://localhost:8081/ws?access_token=${accessToken}`);

        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');

            const message = JSON.parse(messages[0])

            if (message.type === 'insert') {
                const payload = message.payload
                const sender = message.sender
                const obj = {
                    "sender_id": sender.id,
                    "sender_image_url": sender.image_url,
                    "sender_name": sender.name,
                    "content": payload.content,
                    "created_at": payload.created_at
                }

                addMessageToPage(obj)
            }
        };
    };
}

init()

var receiverId = null

const messagesElement = document.getElementById("messages")
function openChat() {
    const id = prompt("Insert user id of user you want to chat with")
    messagesElement.innerHTML = ''
    receiverId = id
    showChat()
}

async function search() {
    const query = document.getElementById("search").value
    if (!query.length) {
        alert("input at least 1 letter")
        return
    }
    messagesElement.innerHTML = ''
    const apiUrl = "http://localhost:8081/search?q=" + query

    const headers = {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${sessionStorage.getItem("accessToken")}`
    }
    try {
        const response = await sendFetchRequest(apiUrl, "GET", headers)
        const data = await response.json()
        data.forEach(d => {
            addMessageToPage(d)
        })
    } catch (error) {

    }
}

async function showChat() {
    const apiUrl = "http://localhost:8081/messages?receiver_id=" + receiverId

    const headers = {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${sessionStorage.getItem("accessToken")}`
    }
    try {
        const response = await sendFetchRequest(apiUrl, "GET", headers)
        const data = await response.json()
        data.forEach(d => {
            addMessageToPage(d)
        })
    } catch (error) {

    }
}

function padTo2Digits(num) {
    return num.toString().padStart(2, '0');
}

function formatDate(date) {
    return (
        [
            date.getFullYear(),
            padTo2Digits(date.getMonth() + 1),
            padTo2Digits(date.getDate()),
        ].join('-') +
        ' ' +
        [
            padTo2Digits(date.getHours()),
            padTo2Digits(date.getMinutes()),
            padTo2Digits(date.getSeconds()),
        ].join(':')
    );
}

async function addMessageToPage(message) {
    const userId = sessionStorage.getItem("userId")
    const element = document.createElement('li')
    element.classList.add('card', 'm-2')

    const leftSide = `<div class="col-sm-2 avatar-container float-end">
        <img src="${message.sender_image_url}"
            alt="" srcset="" >
        <p class="avatar-username">${message.sender_name}</p>
        </div>
        <div class="col-sm-10">
            <p>${message.content}</p>
        </div>`
    const rightSide = `
        <div class="col-sm-10">
            <p>${message.content}</p>
        </div>
        <div class="col-sm-2 avatar-container float-end">
        <img src="${message.sender_image_url}"
            alt="" srcset="" >
        <p class="avatar-username">${message.sender_name}</p>
        </div>`

    element.innerHTML = `
    <li class="card">
        <div class="card-body">
            <div class="row">
                ${userId == message.sender_id ? leftSide : rightSide}
            </div>
            <div class="row">
                <p class="timestamp">${formatDate(new Date(message.created_at))}</p>
            </div>
        </div>
    </li>`
    messagesElement.append(element)
    element.scrollIntoView({ behavior: "smooth" })
}

async function sendFetchRequest(apiUrl, method, headers, payload) {
    try {
        const response = await fetch(apiUrl, {
            method,
            headers,
            ...(!!payload ? { body: JSON.stringify(payload) } : {})
        })
        if (!response.ok) throw new Error(response.statusText)

        return response
    } catch (error) {
        console.error(error.message)
        alert("Maaf terjadi kesalahan")
    }
}

async function sendMessage() {
    if (!receiverId) {
        alert("Press 'Open Chat' first and enter user ID of user you want to chat with")
        return
    }
    const message = document.getElementById("message").value

    const apiUrl = "http://localhost:8081/messages"
    const payload = {
        "content": message,
        "receiver_id": parseInt(receiverId)
    }

    const headers = {
        "Content-Type": "application/json",
        "Authorization": `Bearer ${sessionStorage.getItem("accessToken")}`
    }
    try {
        const response = await sendFetchRequest(apiUrl, "POST", headers, payload)
        await response.json()
    } catch (error) {

    }
}

function logout() {
    sessionStorage.clear();
    location.reload()
}