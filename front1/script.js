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
        const msg = document.createElement('li');
        msg.innerText = JSON.parse(res.data).message;
        messages.appendChild(msg)
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