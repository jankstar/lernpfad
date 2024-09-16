<template>
    <q-dialog ref="dialogRef" @hide="onDialogHide" html>
        <q-card class="q-dialog-plugin">
            <q-card-section>
                <div class="text-h6">{{ title }}</div>
            </q-card-section>

            <q-card-section class="q-pt-none">
                {{ message }}
            </q-card-section>

            <q-card-section>
                <q-input v-model="field1" :label="fieldName1" :placeholder="fieldName1" required autofocus />
                <q-input v-model="field2" :label="fieldName2" :placeholder="fieldName2" required />
            </q-card-section>

            <q-card-actions align="right">
                <q-btn flat icon="cancel" label="Abbruch" v-close-popup />
                <q-btn flat icon="publish" label="Sichern" @click="onOKClick(field1, field2)" />
            </q-card-actions>
        </q-card>
    </q-dialog>
</template>

<script>
import { useDialogPluginComponent } from "quasar";

export default {
    props: {
        title: {
            type: String,
            default: "",
        },
        message: {
            type: String,
            default: "",
        },
        fieldName1: {
            type: String,
            default: "",
        },
        fieldName2: {
            type: String,
            default: "",
        },
    },

    emits: [
        // REQUIRED; need to specify some events that your
        // component will emit through useDialogPluginComponent()
        ...useDialogPluginComponent.emits,
    ],

    data: () => {
        return {
            field1: "",
            field2: "",
        };
    },
    setup() {
        // REQUIRED; must be called inside of setup()
        const { dialogRef, onDialogHide, onDialogOK, onDialogCancel } = useDialogPluginComponent();
        // dialogRef      - Vue ref to be applied to QDialog
        // onDialogHide   - Function to be used as handler for @hide on QDialog
        // onDialogOK     - Function to call to settle dialog with "ok" outcome
        //                    example: onDialogOK() - no payload
        //                    example: onDialogOK({ /*.../* }) - with payload
        // onDialogCancel - Function to call to settle dialog with "cancel" outcome

        return {
            // This is REQUIRED;
            // Need to inject these (from useDialogPluginComponent() call)
            // into the vue scope for the vue html template
            dialogRef,
            onDialogHide,

            // other methods that we used in our vue html template;
            // these are part of our example (so not required)
            onOKClick(iField1, iField2) {
                // on OK, it is REQUIRED to
                // call onDialogOK (with optional payload)
                onDialogOK({ field1: iField1, field2: iField2 });
                // or with payload: onDialogOK({ ... })
                // ...and it will also hide the dialog automatically
            },

            // we can passthrough onDialogCancel directly
            onCancelClick: onDialogCancel,
        };
    },
};
</script>
