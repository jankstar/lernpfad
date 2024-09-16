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
    name: "DiagramPage",
    props: ["user", "server"],
    data: () => {
        return {
            chart: undefined,
            data: {
                nodes: [
                    { name: "Deutsch Klasse 4", depth: 0 },
                    { name: "D402 Abschreibtexte", depth: 1 },
                    { name: "D401 Adjektive", depth: 1 },
                    { name: "D403 Akusativobjekte", depth: 1 },
                    { name: "D404 Aufforderungssätze", depth: 3 },
                    { name: "D405 Satzarten", depth: 2 },
                    { name: "Deutsch Klasse 5/6", depth: 4 },
                ],
                links: [
                    { source: "Deutsch Klasse 4", target: "D402 Abschreibtexte", value: 19 },
                    { source: "Deutsch Klasse 4", target: "D401 Adjektive", value: 21 },
                    { source: "Deutsch Klasse 4", target: "D403 Akusativobjekte", value: 31 },
                    { source: "Deutsch Klasse 4", target: "D405 Satzarten", value: 29 },
                    //
                    { source: "D401 Adjektive", target: "D404 Aufforderungssätze", value: 21 },
                    { source: "D403 Akusativobjekte", target: "D404 Aufforderungssätze", value: 31 },
                    { source: "D405 Satzarten", target: "D404 Aufforderungssätze", value: 29 },
                    //
                    { source: "D402 Abschreibtexte", target: "Deutsch Klasse 5/6", value: 19 },
                    { source: "D404 Aufforderungssätze", target: "Deutsch Klasse 5/6", value: 81 },
                ],
            },
        };
    },

    async created() {},
    async mounted() {
        if (this.chart == undefined) {
            this.chart = echarts.init(document.getElementById("chart"));
        }
        if (!this.chart || this.user.name == "Gast") {
            this.$router.push("/");
            return;
        }
        //this.Chart.showLoading();
        //todo lesen Daten
        //this.Chart.hideLoading();
        this.chart.setOption({
            title: {
                text: "Lernpfad Diagramm Deutsch Klasse 4 - 5/6",
            },
            tooltip: {
                trigger: "item",
                triggerOn: "mousemove",
            },
            series: [
                {
                    type: "sankey",
                    data: this.data.nodes,
                    links: this.data.links,
                    emphasis: {
                        focus: "adjacency",
                    },
                    lineStyle: {
                        color: "gradient",
                        curveness: 0.5,
                    },
                },
            ],
        });
    },
    unmounted() {
        if (this.chart) {
            echarts.dispose(this.chart);
        }
    },

    methods: {},
});
</script>
