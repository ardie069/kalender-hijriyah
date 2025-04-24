<template>
    <div class="mb-4 text-left">
        <label for="method" :class="[darkMode ? 'text-gray-400' : 'text-gray-600', 'text-sm mb-4']">Metode Perhitungan:</label>
        <select id="method" :value="selectedMethod" @change="onMethodChange" :class="[
            methodClass,
            'mt-1 block w-full rounded-lg p-2 pl-2 pr-8 focus:ring-2 focus:outline-none appearance-none',
            darkMode ? 'border-gray-600 focus:ring-gray-500' : 'border-gray-400 focus:ring-gray-300'
        ]">
            <option value="global">Global</option>
            <option value="hisab">Hisab</option>
            <option value="rukyat">Rukyat</option>
        </select>
    </div>
</template>

<script>
export default {
    props: {
        darkMode: Boolean,
        selectedMethod: String
    },
    computed: {
        methodClass() {
            return this.darkMode ? "bg-gray-700 text-white" : "bg-gray-300 text-black";
        }
    },
    methods: {
        onMethodChange(event) {
            this.$emit("update:selectedMethod", event.target.value);
        }
    }
};
</script>
