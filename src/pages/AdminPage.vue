<template>
    <q-page>
        <div>
            <q-tabs v-model="tab" inline-label>
                <q-tab name="kurse" icon="movie" label="Kurse" />
                <q-tab name="user" icon="people" label="Benutzer" />
                <q-tab name="statistik" icon="slow_motion_video" label="Statistik" />
            </q-tabs>
            <q-separator />
            <!-- Kurse -->
            <q-tab-panels v-model="tab" animated>
                <q-tab-panel name="kurse">
                    <!-- Tab Kurs-->
                    <div class="q-pa-md q-gutter-sm row">
                        <q-select v-model="user_id" :options="getSubjects()" label="User" style="width: 100px" @popup-hide="onChangeUser" />
                        <q-select v-model="kurs_id" :options="getKurse()" label="Kurs" style="width: 500px" @popup-hide="onChangeKurs" />

                        <q-btn
                            icon="published_with_changes"
                            aria-label="Prüfen"
                            @click="onCheck()"
                            label="Prüfen"
                            style="background: goldenrod; color: white"
                        >
                            <q-tooltip>Änderungen prüfen</q-tooltip>
                        </q-btn>
                        <q-btn
                            icon="restart_alt"
                            aria-label="Zurücksetzen"
                            @click="onReset()"
                            label="Zurücksetzen"
                            style="background: goldenrod; color: white"
                        >
                            <q-tooltip>Änderungen zurücksetzen und Daten neu laden</q-tooltip>
                        </q-btn>
                        <q-btn
                            icon="file_upload"
                            aria-label="Sichern"
                            @click="onSave()"
                            label="Sichern"
                            style="background: goldenrod; color: white"
                        >
                            <q-tooltip>Daten Sichern - übertragen zum Server</q-tooltip>
                        </q-btn>
                    </div>
                    <q-input v-model="kurs_json" filled type="textarea" input-style="height: 60em" />
                </q-tab-panel>
                <!-- Benutzer und Berechtigungen-->

                <q-tab-panel name="user">
                    <q-input borderless dense debounce="300" v-model="filter" placeholder="Search" outlined style="width: 200px">
                        <template v-slot:append>
                            <q-icon name="search" />
                        </template>
                    </q-input>
                    <div class="text-h6">Gruppe</div>
                    <q-btn @click="onAddGruppe()" label="Neu" style="background: goldenrod; color: white; margin-top: 5px" />
                    <q-btn @click="onRemoveGruppe()" label="Löschen" style="background: goldenrod; color: white; margin-top: 5px" />
                    <div class="q-pa-md">
                        <q-table
                            :rows="user.admin.allSPolicy"
                            row-key="id"
                            :filter="filter"
                            v-model:pagination="pagination1"
                            :rows-per-page-options="[0]"
                            selection="single"
                            v-model:selected="selected1"
                            style="height: 300px"
                        />
                    </div>
                    <div class="text-h6">Benutzer</div>
                    <q-btn @click="onAddUser()" label="Neu" style="background: goldenrod; color: white; margin-top: 5px" />
                    <q-btn @click="onRemoveUser()" label="Löschen" style="background: goldenrod; color: white; margin-top: 5px" />
                    <div class="q-pa-md">
                        <q-table
                            :rows="user.admin.allGroups"
                            row-key="id"
                            :filter="filter"
                            v-model:pagination="pagination2"
                            :rows-per-page-options="[0]"
                            selection="single"
                            v-model:selected="selected2"
                            style="height: 300px"
                        />
                    </div>
                </q-tab-panel>
            </q-tab-panels>
        </div>
    </q-page>
</template>

<script>
import { Kurs } from "src/lib/master_data";
import { defineComponent } from "vue";
import addDialog from "src/components/addDialog.vue";

export default defineComponent({
    name: "AdminPage",
    props: ["user", "server"],
    data: () => {
        return {
            tab: "user",
            user_id: { value: "Gast", label: "Gast" },
            kurse: [],
            kurs_id: undefined,
            kurs_json: "",
            filter: "",
            pagination1: "",
            pagination2: "",
            selected1: [],
            selected2: [],
        };
    },

    async created() {
        if (!this.user.role.includes("admin")) {
            this.$router.push("/");
        }
        this.onChangeUser();
    },

    methods: {
        async onSave() {
            this.onCheck(true);
        },
        async onReset() {
            let lUser_id = this.user_id;
            let lKurs_id = this.kurs_id;
            await this.onChangeUser();
            this.user_id = lUser_id;
            this.kurs_id = lKurs_id;
            this.onChangeKurs();
        },
        onCheck(iSave) {
            let lSave = iSave || false;
            if (this.user_id && this.user_id.value && this.kurs_id && this.kurs_id.value) {
                for (let i = 0; i < this.kurse.length; i += 1) {
                    if (this.kurse[i].id == this.kurs_id.value) {
                        let lObj_json = this.kurs_json.replaceAll("<br>\n", "<br>");
                        let lObj = JSON.parse(lObj_json);
                        this.kurse[i] = new Kurs(lObj);
                        if (lSave) {
                            let that = this;
                            this.user.setItem("Kurs_" + this.kurse[i].id, this.user_id.value, this.kurse[i]).then((res) => {
                                if (res.status != 200) {
                                    let lMessage = res.message || "Fehler";
                                    that.$q.notify({
                                        message: `Fehler beim Sichern:'${lMessage}'`,
                                        color: "negative",
                                        icon: "warning",
                                    });
                                } else {
                                    that.$q.notify({
                                        message: `Daten gesichert - Daten sind erst nach neuem Login sichtbar.`,
                                        type: "positive",
                                        time: 500,
                                    });
                                }
                            });
                        }
                        this.kurs_json = JSON.stringify(this.kurse[i], null, 5);
                        this.kurs_json = this.kurs_json.replaceAll("<br>", "<br>\n");
                        break;
                    }
                }
            }
        },
        async onChangeUser() {
            this.kurs_id = undefined;
            this.kurs_json = "";
            if (this.user_id && this.user_id.value) {
                this.kurse = await this.user.getKurseByUser(this.user_id.value);
            } else {
                this.kurse = [];
            }
        },

        onChangeKurs() {
            if (this.kurs_id && this.kurs_id.value && this.user_id && this.user_id.value) {
                for (let k of this.kurse) {
                    if (k.id == this.kurs_id.value) {
                        this.kurs_json = JSON.stringify(k, null, 5);
                        this.kurs_json = this.kurs_json.replaceAll("<br>", "<br>\n");
                        break;
                    }
                }
            } else {
                this.kurs_json = "";
            }
        },
        getKurse() {
            let lRet = [];
            if (this.kurse && this.kurse.length > 0) {
                for (let u of this.kurse) {
                    lRet.push({ value: u.id, label: u.id + " - " + u.main_title });
                }
            }
            return lRet;
        },
        getSubjects() {
            let lRet = [];
            if (this.user.admin && this.user.admin.allGroups && this.user.admin.allGroups.length > 0) {
                for (let u of this.user.admin.allGroups) {
                    let lEle = {};
                    lEle.value = u.User;
                    lEle.label = u.User;
                    if (lRet.findIndex((element) => element.value == lEle.value) < 0) {
                        lRet.push(lEle);
                    }
                }
            }
            return lRet;
        },

        onAddGruppe() {
            this.$q
                .dialog({
                    component: addDialog,
                    componentProps: {
                        title: "Neu",
                        message: `Erfassen einer neuen Gruppe:`,
                        fieldName1: "Gruppe:",
                        fieldName2: "Kurs:",
                    },
                })
                .onOk((data) => {
                    let that = this;
                    this.user.addGruppe(data.field1, data.field2).then((data) => {
                        if (data && data.status && data.status == 200) {
                            this.$q.notify({
                                message: `Daten gesichert.`,
                                type: "positive",
                                time: 500,
                            });
                            that.user.getRole();
                        } else {
                            if (!data.message) {
                                data.message = "Sichern nicht möglich.";
                            }
                            this.$q.notify({
                                message: `Fehler: ${data.message}`,
                                color: "negative",
                                icon: "warning",
                            });
                        }
                    });
                });
        },

        onRemoveGruppe() {
            if (!this.selected1 || this.selected1.length == 0) {
                this.$q.notify({
                    message: "Bitte eine Zeile aus der Tabelle Gruppe auswählen.",
                    color: "negative",
                    icon: "warning",
                });
                return;
            }
            this.$q
                .dialog({
                    title: "Frage <span class='material-icons'>question_answer</span>",
                    message: `Soll der Eintrag <br/> - Gruppe: ${this.selected1[0].Gruppe} <br/> - Kurs: ${this.selected1[0].Kurs} <br/>gelöscht werden?`,
                    cancel: { label: "Nein", flat: true, "v-close-popup": true, icon: "clear" },
                    ok: { label: "Ja", flat: true, icon: "check" },
                    persistent: true,
                    html: true,
                })
                .onOk(() => {
                    let that = this;
                    this.user.removeGruppe(this.selected1[0].Gruppe, this.selected1[0].Kurs).then((data) => {
                        if (data && data.status && data.status == 200) {
                            this.$q.notify({
                                message: `Daten gesichert.`,
                                type: "positive",
                                time: 500,
                            });
                            that.user.getRole();
                            that.selected1 = [];
                        } else {
                            if (!data.message) {
                                data.message = "Sichern nicht möglich.";
                            }
                            this.$q.notify({
                                message: `Fehler: ${data.message}`,
                                color: "negative",
                                icon: "warning",
                            });
                        }
                    });
                });
        },

        onAddUser() {
            this.$q
                .dialog({
                    component: addDialog,
                    componentProps: {
                        title: "Neu",
                        message: `Erfassen eines neuen Benutzers:`,
                        fieldName1: "User:",
                        fieldName2: "Gruppe:",
                    },
                })
                .onOk((data) => {
                    let that = this;
                    this.user.addUser(data.field1, data.field2).then((data) => {
                        if (data && data.status && data.status == 200) {
                            this.$q.notify({
                                message: `Daten gesichert.`,
                                type: "positive",
                                time: 500,
                            });
                            that.user.getRole();
                        } else {
                            if (!data.message) {
                                data.message = "Sichern nicht möglich.";
                            }
                            this.$q.notify({
                                message: `Fehler: ${data.message}`,
                                color: "negative",
                                icon: "warning",
                            });
                        }
                    });
                });
        },
        onRemoveUser() {
            if (!this.selected2 || this.selected2.length == 0) {
                this.$q.notify({
                    message: "Bitte eine Zeile aus der Tabelle Benutzer auswählen.",
                    color: "negative",
                    icon: "warning",
                });
                return;
            }
            this.$q
                .dialog({
                    title: "Frage <span class='material-icons'>question_answer</span>",
                    message: `Soll der Eintrag <br/> - User: ${this.selected2[0].User} <br/> - Gruppe: ${this.selected2[0].Gruppe} <br/>gelöscht werden?`,
                    cancel: { label: "Nein", flat: true, "v-close-popup": true, icon: "clear" },
                    ok: { label: "Ja", flat: true, icon: "check" },
                    persistent: true,
                    html: true,
                })
                .onOk(() => {
                    let that = this;
                    this.user.removeUser(this.selected2[0].User, this.selected2[0].Gruppe).then((data) => {
                        if (data && data.status && data.status == 200) {
                            this.$q.notify({
                                message: `Daten gesichert`,
                                type: "positive",
                                time: 500,
                            });
                            that.user.getRole();
                            that.selected2 = [];
                        } else {
                            if (!data.message) {
                                data.message = "Sichern nicht möglich.";
                            }
                            this.$q.notify({
                                message: `Fehler: ${data.message}`,
                                color: "negative",
                                icon: "warning",
                            });
                        }
                    });
                });
        },
    },
});
</script>
