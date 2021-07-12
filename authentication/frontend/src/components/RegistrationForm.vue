<template>
    <form id="registrationForm" action="#" method="POST" @submit.prevent="register()">
        <article class="input">
            <input v-model="email" v-on:input="(e) => {
                taken = taken.filter(t => {
                    return t != 'email'
                })
            }" name="email" type="email" autocomplete="email" id="email" placeholder="Email" required>
            <label for="Email" v-bind:class="taken.includes('email') ? 'taken' : ''">{{taken.includes('email') ? "Email taken" : "Email"}}</label>
        </article>

        <article class="input">
            <input v-model="username" v-on:input="(e) => {
                taken = taken.filter(t => {
                    return t != 'username'
                })
            }" name="username" autocomplete="username" id="username" placeholder="Username" required>
            <label for="username" v-bind:class="taken.includes('username') ? 'taken' : ''">{{taken.includes('username') ? "Username taken" : "Username"}}</label>
        </article>
        <article class="input">
            <input v-model="password" name="password" autocomplete="current-password" id="password" type="password" placeholder="Password" required>
            <label for="password">Password</label>
        </article>
        <input id="registerBtn" type="submit" v-bind:class="{
            success: success
        }" v-bind:value="registering ? 'Registering...' : success ? 'Registration complete.' : 'Register'" v-bind:disabled="registering || success">
    </form>
</template>

<script>
import config from '../config'
export default {
  name: 'RegistrationForm',
  data: function() {
    return {
      loginFailed: false,
      loginFailReason: null,
      email: "",
      username: "",
      password: "",
      registering: false,
      success: false,
      taken: []
    }
  },
  computed: {
        requiredEntries: function() {
            let required = [];

            for (let field of this.fields) {
                if (this.fields[field].required && this.fields[field].value == "") {
                    required.push(field)
                }
            }

            return required
        }
  },
  methods: {
      register: function() {
          this.registering = true;

          const data = JSON.stringify({"username":this.username,"password":this.password,"email": this.email});

          console.log(data)

          var xhr = new XMLHttpRequest();
          xhr.withCredentials = true;

            const _this = this;
          xhr.addEventListener("readystatechange", function() {
              if(this.readyState === 4) {
                    const resp = this.responseText

                    _this.registering = false;

                    console.log(resp)

                    if (this.status != 200) {
                        console.log(resp == "username taken")
                        switch(resp) {
                            case "username taken":
                                _this.taken.push("username")
                                break;
                        }

                        return
                    }

                    _this.success = true

                    setTimeout(() => {
                        document.querySelector('.page.showing').classList.remove('showing')
                        document.querySelector('.page.login').classList.add('showing')
                    }, 1000)
              }
          });

          xhr.open("POST", `${config.API_URL}/register`);
          xhr.setRequestHeader("Content-Type", "application/json");

          xhr.send(data);
      }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped>
    @import 'styles.scss';
</style>
