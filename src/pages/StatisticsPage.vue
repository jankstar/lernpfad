<template>
    <q-page>
        <div class="row">
            <q-card style="height: 550px; width: 100%">
                <q-card-section>
                    <div id="chart" style="height: 520px"></div>
                </q-card-section>
            </q-card>
        </div>
    </q-page>
</template>
<style lang="sass" scoped>
.my-card
  width: 100%
  max-width: 300px
  height: 420px
</style>
<script>
import { defineComponent } from "vue";
import * as echarts from "echarts";

export default defineComponent({
    name: "StatisticsPage",
    props: ["user", "server"],
    data: () => {
        return {
            chart: undefined,
        };
    },

    created() {},
    mounted() {
        this.chart = echarts.init(document.getElementById("chart"));
        if (this.user.name == "Gast" || this.user.kurse.length == 0) {
            this.$router.push("/");
            return;
        }
        //this.Chart.showLoading();
        //todo lesen Daten
        //this.Chart.hideLoading();

        let position = [
            { left: "0%", right: "90%" },
            { left: "10%", right: "80%" },
            { left: "20%", right: "70%" },
            { left: "30%", right: "60%" },
            { left: "40%", right: "50%" },
            { left: "50%", right: "40%" },
            { left: "60%", right: "30%" },
            { left: "70%", right: "20%" },
            { left: "80%", right: "10%" },
            { left: "90%", right: "0%" },
        ];

        let lTitle = [
            {
                text: "Meine Lernpfade",
                left: "center",
            },
        ];
        let lSeries = [];
        let lYPos = 0;
        for (let k = 0; k < this.user.kurse.length; k += 1) {
            if (this.user.kurse[k].start != "") {
                lYPos += 1;
                let lPoints = this.user.kurse[k].getPoints();
                let lMaxPoints = this.user.kurse[k].getMaxPoints();
                lTitle.push({
                    subtext: this.user.kurse[k].main_title.substring(0, 40) + `\nPunkte: ${lPoints} von ${lMaxPoints}`,
                    left: position[0].left,
                    top: position[lYPos].left,
                    textAlign: "left",
                });
                for (let t = 0; t < this.user.kurse[k].tasks.length; t += 1) {
                    let eOffen = 0;
                    let eRichtig = 0;
                    if (this.user.kurse[k].tasks[t].icon == "quiz") {
                        for (let question of this.user.kurse[k].tasks[t].questions) {
                            eOffen += question.getMaxPoints();
                            if (this.user.kurse[k].tasks[t].check != "") {
                                eRichtig += question.getPoints();
                                eOffen -= question.getPoints();
                            }
                        }
                    } else {
                        //hier ist alles andere außer Quiz
                        if (this.user.kurse[k].tasks[t].start == "") {
                            eOffen = 100;
                        } else {
                            eRichtig = 100;
                        }
                    }
                    lSeries.push({
                        type: "pie",
                        radius: "50%",
                        center: ["50%", "50%"],
                        data: [
                            {
                                name: "richtig/fertig",
                                label: { show: false },
                                itemStyle: { color: "#00b300" },
                                value: eRichtig,
                            },
                            {
                                name: "falsch/offen",
                                label: { show: false },
                                itemStyle: { color: "#0000b3" },
                                value: eOffen,
                            },
                        ],
                        left: position[t + 1].left,
                        right: position[t + 1].right,
                        top: position[lYPos].left,
                        bottom: position[lYPos].right,
                    });
                }
            }
        }

        let option = {
            legend: {
                top: "5%",
                left: "center",
            },
            title: lTitle,
            series: lSeries,
        };

        this.chart.setOption(option);
    },
    unmounted() {
        if (this.chart) {
            echarts.dispose(this.chart);
        }
    },

    methods: {},
});
</script>
