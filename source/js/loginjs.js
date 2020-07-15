var app = new Vue({
  el: '#loginModal',
  methods:{
    closeModal: function () {
      $('#loginModal').modal('hide');
    },
    checkPublicKey: function(publickey) {
      console.log("Checking Public Key: ", publickey)
      axios.get(`${app.baseUrl}/account/${publickey}`).then(
        function(response) {
          console.log(response.data)
          if (response.data.isRoot === undefined) {
            app.closeModal()
            alert("Public Key has no root access");
          }
          if (response.data.isRoot === true) {
            this.$cookies.set("publicKey", publickey, 60 * 60 * 24 * 30)
            app.closeModal()
            alert("Successfully Verified! The browser will now refresh.");

            location.reload();
          }
        }
      )
    }
  },
  data: {
    baseUrl: 'http://localhost:8010',
    accountinfo: {
      publickey: '',
    },
  },
  delimiters: ["<%","%>"]
})
