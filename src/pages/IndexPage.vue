<template>
    <q-page>
        <div>
            <div>
                <q-card>
                    <q-card-section class="q-pa-md row">
                        <div class="row">
                            <div style="margin-top: 20px">Filter:</div>
                            &nbsp; &nbsp;
                            <q-select
                                outlined
                                clearable
                                v-model="setting.klasse"
                                :options="user.grade_cat"
                                label="Klasse"
                                style="min-width: 150px"
                            />
                            <q-select
                                outlined
                                clearable
                                v-model="setting.fach"
                                :options="user.subject_cat"
                                label="Fach"
                                style="min-width: 150px"
                            />
                            <q-select
                                outlined
                                v-model="setting.sel_status"
                                :options="filter_status"
                                label="Status"
                                style="min-width: 150px"
                            />
                        </div>
                        &nbsp; &nbsp;
                        <div class="row">
                            <div style="margin-top: 20px">Suche:</div>
                            &nbsp; &nbsp;
                            <q-input outlined clearable v-model="setting.suche_text" label="Begriff" style="min-width: 200px" />
                        </div>
                        &nbsp; &nbsp;
                        <div class="row">
                            <div style="margin-top: 20px">Sortieren:</div>
                            &nbsp; &nbsp;
                            <q-select
                                outlined
                                v-model="setting.sortierung"
                                :options="['', 'Klasse', 'Fach', 'Start']"
                                label="Merkmal"
                                style="min-width: 150px"
                                @popup-hide="onChangeSort"
                            />
                        </div>
                    </q-card-section>
                </q-card>
            </div>
            <div class="flex flex-center">
                <div v-if="user.kurse.length != 0" class="q-pa-md row items-start q-gutter-md">
                    <template v-for="kurs in user.kurse" v-bind:key="kurs.id">
                        <q-card v-if="user.kurse.length != 0 && onDisplay(kurs) == true" class="my-card">
                            <q-btn style="border: 0; margin: 0; width: 100%" stretch flat @click="onGoKurs(kurs.id)">
                                <img :src="kurs.image" style="object-fit: cover; width: 100%; height: 180px" flat />
                                <span v-if="kurs.start" class="text-h4" style="position: absolute; left: 10%; top: 10%">
                                    <span class="material-icons" style="color: #a3005f"> verified </span>
                                </span>
                                <span
                                    v-if="!kurs.start && isNew(kurs.publication)"
                                    class="text-h4"
                                    style="position: absolute; left: 10%; top: 10%"
                                >
                                    <span class="q-pt-none" style="color: #a3005f"> neu </span>
                                </span>
                            </q-btn>

                            <q-card-section>
                                <div class="text-h6">{{ kurs.main_title }}</div>
                                <div class="text-subtitle2">{{ kurs.sub_title }}</div>
                            </q-card-section>

                            <q-card-section class="q-pt-none">
                                {{ kurs.description_short }}
                            </q-card-section>
                            <q-card-section class="q-pt-none">
                                Klasse: {{ kurs.grade }}&nbsp; &nbsp; Fach: {{ kurs.subject }}<br />
                                Teilnehmer: {{ kurs.enrollments }}
                                &nbsp; &nbsp; Score:
                                <q-rating
                                    v-model="kurs.rating"
                                    size="01.0em"
                                    color="yellow"
                                    icon="star_border"
                                    icon-selected="star"
                                    readonly
                                />
                            </q-card-section>
                        </q-card>
                    </template>
                </div>
            </div>
        </div>
    </q-page>
</template>
<style lang="sass" scoped>
.my-card
  width: 100%
  max-width: 300px
  height: 450px
</style>
<script>
import { defineComponent } from "vue";

export default defineComponent({
    name: "IndexPage",
    props: ["user", "server"],
    data: () => {
        return {
            setting: {},
            filter_status: [
                { value: "1", label: "alle Kurse" },
                { value: "2", label: "meine laufenden Kurse" },
                { value: "3", label: "meine beendeten Kurse" },
                { value: "4", label: "offene Kurse" },
                { value: "5", label: "offene und laufende Kurse" },
            ],
        };
    },

    async created() {
        this.setting = this.user.setting;
    },

    methods: {
        onDisplay(iKurs) {
            this.user.setSetting(this.setting);

            if (this.setting.sel_status.value == "2" && (iKurs.start == "" || iKurs.terminated != "")) {
                // "meine laufenden Kurse"
                return false;
            }
            if (this.setting.sel_status.value == "3" && iKurs.terminated == "") {
                // "meine beendeten Kurse"
                return false;
            }
            if (this.setting.sel_status.value == "4" && iKurs.start != "") {
                // "offene Kurse"
                return false;
            }

            if (this.setting.sel_status.value == "5" && iKurs.terminated != "") {
                // "offene und laufende Kurse"
                return false;
            }

            if (this.setting.klasse && iKurs.grade != this.setting.klasse) {
                return false;
            }

            if (this.setting.fach && iKurs.subject != this.setting.fach) {
                return false;
            }

            if (!this.setting.suche_text || this.setting.suche_text == "") {
                return true;
            }
            if (iKurs.main_title && iKurs.main_title.toUpperCase().includes(this.setting.suche_text.toUpperCase())) {
                return true;
            }
            if (iKurs.sub_title && iKurs.sub_title.toUpperCase().includes(this.setting.suche_text.toUpperCase())) {
                return true;
            }
            if (iKurs.description_short && iKurs.description_short.toUpperCase().includes(this.setting.suche_text.toUpperCase())) {
                return true;
            }
            if (iKurs.description && iKurs.description.toUpperCase().includes(this.setting.suche_text.toUpperCase())) {
                return true;
            }
            return false;
        },

        onChangeSort() {
            this.user.sortKurseByFeature(this.setting.sortierung);
        },
        isNew(iDatum) {
            // 10 Tage lang ist der Kurs als "NEU" markiert
            let iDateTime = iDatum.split(" ");
            if (!iDatum || iDateTime.length < 2) {
                return false;
            }
            let iDateArray = iDateTime[0].split(".");
            if (iDateArray.length < 2) {
                return false;
            }
            let lYear = +iDateArray[2];
            let lMon = +iDateArray[1] - 1;
            let lDay = +iDateArray[0];
            let lDatumTest = new Date(lYear, lMon, lDay);
            let lDatumNow = new Date();
            let lDiff = Math.floor((lDatumNow - lDatumTest) / (1000 * 60 * 60 * 24));
            if (lDiff > 10 || lDiff < -10) {
                return false;
            }
            return true;
        },
        onGoKurs(kurs_id) {
            this.$router.push("/kurs/" + kurs_id);
        },
    },
});
</script>
