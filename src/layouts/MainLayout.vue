<template>
  <q-layout view="lHh Lpr lFf" class="tw-font-sans">
    <q-header elevated>
      <q-toolbar>
        <q-btn flat dense round icon="menu" aria-label="Menü" @click="toggleLeftDrawer">
          <q-tooltip>Menü</q-tooltip>
        </q-btn>
        <q-btn flat dense round icon="menu_book" aria-label="Home" @click="$router.push('/')">
          <q-tooltip>Startseite</q-tooltip>
        </q-btn>
        <q-btn v-if="user.name != 'Gast'" flat dense round icon="mediation" aria-label="Diagramm"
          @click="$router.push('/diagram')">
          <q-tooltip>Diagramm mit den Kursfolgen</q-tooltip>
        </q-btn>
        <q-btn v-if="user.name != 'Gast'" flat dense round icon="slow_motion_video" aria-label="Statistik"
          @click="$router.push('/statistics')">
          <q-tooltip>Meine Statistik</q-tooltip>
        </q-btn>
        <q-btn v-if="user.role.includes('admin')" flat dense round icon="settings_suggest" aria-label="Administration"
          @click="$router.push('/admin')">
          <q-tooltip>Administration</q-tooltip>
        </q-btn>
        <q-toolbar-title> Lernpfad </q-toolbar-title>

        <div>{{ user.name }}&nbsp; &nbsp;</div>
        <q-btn v-if="user.name == '' || user.name == 'Gast'" flat dense round icon="login" aria-label="Benutzer"
          @click="onLogout()"><q-tooltip>Login</q-tooltip></q-btn>
        <q-btn v-else flat dense round icon="logout" aria-label="Logout"
          @click="onLogout()"><q-tooltip>Logout</q-tooltip></q-btn>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered>
      <q-list>
        <q-item-label header> Wichtige Absprünge </q-item-label>

        <EssentialLink v-for="link in essentialLinks" :key="link.title" v-bind="link" />
      </q-list>
    </q-drawer>

    <q-page-container>
      <q-dialog v-model="loginDialog">
        <q-card style="width: 700px; max-width: 80vw">
          <q-card-section>
            <div class="text-h6">Login</div>
          </q-card-section>

          <q-card-section class="q-pt-none">
            <q-form method="POST" action="/login" autocorrect="off" autocapitalize="off" autocomplete="off"
              spellcheck="false">
              <q-input v-model="loginDialogName" type="text" label="Name" placeholder="Name zur Anzeige" name="uname"
                required autofocus @keydown.enter="onLogin()">
              </q-input>
              <!--q-input v-model="email" filled type="email" placeholder="EMail" name="email">
                            </q-input-->
              <q-separator spaced inset></q-separator>
              <q-input v-model="loginDialogPassword" type="password" label="Passwort" placeholder="Passwort"
                name="password" required @keydown.enter="onLogin()">
              </q-input>
              <q-separator spaced inset></q-separator>
              <q-btn type="button" :disable="!isValid" style="background: goldenrod; color: white" @click="onLogin()"
                @keydown.enter="onLogin()">Login</q-btn>
            </q-form>
          </q-card-section>
        </q-card>
      </q-dialog>
      <router-view v-model.sync="loginDialog" :user="user" :server="server" />
    </q-page-container>
  </q-layout>
</template>

<script>
import { defineComponent } from "vue";
import EssentialLink from "components/EssentialLink.vue";
import { User } from "../lib/master_data";

const linksList = [
  {
    title: "Lernpfad",
    caption: "Lernvideos - für alle, die es wissen wollen",
    icon: "menu_book",
    link: "https://lernpfad.azurewebsites.net/",
  },
  {
    title: "Impressum",
    icon: "admin_panel_settings",
    link: "https://jankstar.de/impressum/",
  },
];

export default defineComponent({
  name: "MainLayout",

  components: {
    EssentialLink,
  },
  data: () => {
    return {
      essentialLinks: linksList,
      leftDrawerOpen: false,
      server: location.protocol + "//" + location.host, //"http://localhost:8081",
      loginDialog: false,
      loginDialogName: "",
      loginDialogPassword: "",
      user: undefined,
    };
  },

  computed: {
    isValid() {
      return this.loginDialogName.length >= 1 && this.loginDialogPassword.length >= 1;
    },
  },

  async created() {
    this.user = new User({ id: "Gast", name: "Gast", server: this.server });
    await this.user.login(this.server, "Gast", "Gast");
    await this.user.loadKurse();
    this.loginDialog = false;
    this.loginDialogName = "";
    this.loginDialogPassword = "";
    this.leftDrawerOpen = false;
  },

  methods: {
    toggleLeftDrawer() {
      this.leftDrawerOpen = !this.leftDrawerOpen;
    },
    onLogin() {
      let lUser = this.loginDialogName;
      let lPassword = this.loginDialogPassword;
      this.loginDialog = false;
      this.loginDialogName = "";
      this.loginDialogPassword = "";
      this.user = new User({ id: lUser, name: lUser });

      let that = this;
      this.user.login(this.server, lUser, lPassword).then((res) => {
        if (res.status == 200) {
          that.user.loadKurse();
        } else {
          //wieder auf Gast setzen
          that.user = new User();
          that.user.login(this.server, "Gast", "Gast").then(() => {
            that.user.loadKurse();
          });

          if (res.message) {
            this.$q.notify({
              message: res.message,
              color: "negative",
              icon: "warning",
            });
          }
        }
      });
    },

    onLogout() {
      this.$router.push("/");
      this.loginDialogName = "";
      this.loginDialogPassword = "";
      if (!this.user || this.user.name != "Gast") {
        this.user.logout();
        this.user = new User("Gast", "Gast");
        let that = this;
        this.user.login(this.server, "Gast", "Gast").then((res) => {
          if (res.status == 200) {
            that.user.loadKurse();
          } else {
            that.kurse = [];
          }
        });
      } else {
        this.loginDialog = true;
      }
    },
  },
});
</script>
