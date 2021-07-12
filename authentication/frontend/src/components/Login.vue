<template>
    <section>
        <section id="signin" v-bind:class="{ failed: loginFailed, container: true }">
            <article id="body" class="page login" v-bind:class="{ showing: (showingPage == 'login') }">
                <h2 id="title">Login to <span class="serviceName">PRAUXY</span></h2>
                <p id="explanation">Please sign in to your PRAUXY account. This will let you seamlessly connect to PRAUXY-protected services without much hassle</p>
                <LoginForm></LoginForm>
                <hr>
                <a id="transferRegister" href="#" @click="() => { showingPage = 'register'; }">or, create a new account</a>
            </article>
            <article id="body" class="page register" v-bind:class="{ showing: (showingPage == 'register') }">
                <h2 id="title">Register to <span class="serviceName">PRAUXY</span></h2>
                <p id="explanation">By registering, you agree to our terms of service.</p>
                <RegistrationForm></RegistrationForm>
                <hr>
                <a id="transferLogin" href="#" @click="() => { showingPage = 'login'; }">or, login</a>
            </article>
        </section>
        <a href="https://prauxy.app"><img id="logo" src="@/assets/textlogo.svg"></a>
    </section>
</template>

<script>
// import config from '../config'
import LoginForm from './LoginForm.vue'
import RegistrationForm from './RegistrationForm.vue'
export default {
  name: 'Login',
  components: {
      LoginForm,
      RegistrationForm
  },
  data: function() {
    return {
        showingPage: "login"
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
        #badlogin {
            display: none;
            color: rgb(240, 24, 24);
        }
        #badlogin.showing {
            display: inherit;
        }
        .container {
            border-radius: 15px;
            border: 2px solid rgb(236, 236, 236);
            overflow: hidden;
            width: 35vw;
            min-width: 350px;
        }

        /* Portrait */
        @media only screen 
            and (-webkit-min-device-pixel-ratio: 2)
            and (max-device-width: 480px) {
                .container {
                    width: 100%;
                    padding: 15px;
                }
        }

        p {
            padding: 0;
            margin: 0;
        }

        .container.failed {
            border: 2px solid rgb(240, 24, 24);
        }

        #madeBy {
            background-color: rgb(245, 245, 245);
            padding: 15px;
        }

        #body {
            padding: 35px;
            display: flex;
            flex-direction: column;
            align-items: center;
        }

        .page {
            display: none !important;
        }
        .page.showing {
            display: inherit !important;
        }

        a:hover img#logo {
            -webkit-filter: grayscale(0%) opacity(1.0);
            filter: grayscale(0%) opacity(1.0);
            transition: .2s;
        }

        img#logo {
            transition: .2s;
            -webkit-filter: grayscale(100%) opacity(0.5);
            filter: grayscale(100%) opacity(0.5);
            width: 150px;
            padding-top: 15px;
        }

        hr {
            margin-top: 25px;
            margin-bottom: 25px;
            border: 1px solid rgb(218, 218, 218);
        }

        #body a {
            text-decoration: none;
            color: rgb(71, 71, 253);
        }

        #consent {
            display: none;
        }

        .serviceName {
            color: #4B77F0;
        }
</style>
