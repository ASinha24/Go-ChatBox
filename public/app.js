new Vue({
    e1: '#app',

    data: {
        ws: null, //our socket
        newMsg: '', //hold new msg to be sent to the server
        chatcontent: '', // a running list of the chat content displayed on the screen
        email: null,   // email list used for grabbing an avatar
        username: null, //our username
        joined: false // true if email and username have filled
    },

created: function(){
    var self =this,
    this.ws = new WebSocket('ws://'+ window.location.host + '/ws');
    this.ws.addEventListener('message', function(e){
    var msg = JSON.parse(e.data);
    self.chatcontent += '<div class="chip">'
            + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
            + msg.username
            + '</div>'
            + emojione.toImage(msg.message) + '<br/>'; // Parse emojis

    var element = document.getElementById('chat-messages');
    element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
    });
},
methods: {
    send: function () {
        if (this.newMsg != '') {
            this.ws.send(
                JSON.stringify({
                    email: this.email,
                    username: this.username,
                    message: $('<p>').html(this.newMsg).text() // Strip out html
                }
            ));
            this.newMsg = ''; // Reset newMsg
        }
    },
    join: function () {
        if (!this.email) {
            Materialize.toast('You must enter an email', 2000);
            return
        }
        if (!this.username) {
            Materialize.toast('You must choose a username', 2000);
            return
        }
        this.email = $('<p>').html(this.email).text();
        this.username = $('<p>').html(this.username).text();
        this.joined = true;
    },
    gravatarURL: function(email) {
        return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
    }
}
});
