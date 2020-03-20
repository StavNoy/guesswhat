document.addEventListener('DOMContentLoaded', () => {
    const
        input = document.getElementById('input'),
        send = document.getElementById('send'),
        messages = document.getElementById('messages'),
        socket = new WebSocket("ws://localhost:8080/")
        state = document.getElementById('state')
    ;

    socket.onopen = () => {
        state.innerText = 'Connected';
        state.classList.add('connected')
    };

    socket.onmessage = res => {
        [].concat(JSON.parse(res.data)).forEach(msg => {
            const msgEl = document.createElement('li');
            msgEl.innerText = msg.message;
            messages.appendChild(msgEl)
        });
    };

    socket.onclose = socket.onerror = (ev) => {
        state.innerText = 'Disconnected';
        state.classList.remove('connected')
    };

    send.onclick = () => {
        let toSend;

        if (input.value.startsWith('/name ')) {
            toSend = {
                type: 'nickname',
                nickname: input.value.substring('/name '.length)
            }
        } else {
            toSend = {
                type: 'message',
                message: input.value
            }
        }

        socket.send(JSON.stringify(toSend));
    };

    send.focus()
});