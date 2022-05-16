<template>
  <v-dialog v-model="open" persistent max-width="450px" @keydown.esc="open = false" >
    <v-card class="mb-0 pb-0">
      <v-card-title>
        <h2 style="font-size: 18px;font-weight: 500;text-align: center;margin-bottom: 20px;">Время задержки отправления</h2>
      </v-card-title>
      <v-card-text class="mb-0 pb-0">
        <v-form
            ref="formDelayTime"
            v-model="valid"
            lazy-validation
        >
          <div class="time-delay">
            <div class="time-delay-item">
<!--              <h3>Месяцы</h3>-->
              <h3>Недели</h3>
              <h3>Дни</h3>
<!--              <v-text-field dense v-model.number.trim="months" outlined class="mt-2" :rules="[rules.max1]"></v-text-field>-->
              <v-text-field dense v-model.number.trim="weeks" outlined class="mt-2" :rules="[rules.max8, rules.positiveNumber]"></v-text-field>
              <v-text-field dense v-model.number.trim="days" outlined class="mt-2" :rules="[rules.max30, rules.positiveNumber]"></v-text-field>
            </div>
            <div class="time-delay-item">
              <h3>Часы</h3>
              <h3>Минуты</h3>
              <v-text-field dense v-model.number.trim="hours" outlined class="mt-2" :rules="[rules.max24, rules.positiveNumber]"></v-text-field>
              <v-text-field dense v-model.number.trim="minutes" outlined class="mt-2" :rules="[rules.max60, rules.positiveNumber]"></v-text-field>
            </div>
          </div>
        </v-form>

      </v-card-text>
      <v-card-actions class="mt-0 pt-0 pb-3">
        <v-spacer></v-spacer>
        <v-btn color="secondary" dark elevation="1" small outlined @click="$emit('close')" class="mr-3">Отмена</v-btn>
        <v-btn color="teal darken-3" dark elevation="1" small @click="submitTime">Сохранить</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>

export default {
  data: () => ({
    weeks: 0,
    days: 0,
    hours: 0,
    minutes: 0,
    valid: true,
    rules: {
      max24: value => value <= 24 || 'Max 24',
      max60: value => value <= 60 || 'Max 60',
      max30: value => value <= 30 || 'Max 30',
      max8: value => value <= 4 || 'Max 4',
      max1: value => value <= 1 || 'Max 1',
      positiveNumber: value => value >=0 || '> 0',
      box: value => {
        const pattern = /^[a-zA-Z0-9.]+$/
        return pattern.test(value) || 'Только разрешенные символы a-zA-Z0-9..'
      },
    },
  }),
  props: ['open','time'],
  async mounted() {

  },
  watch: {
    open() {

      let _time = this.time / 1000000;

      let d = this.$moment.duration(_time)
      this.months = d.months()
      this.weeks = d.subtract(this.months, 'm').weeks() + this.months*4
      this.days = d.subtract(this.weeks, 'w').days()
      this.hours = d.subtract(this.days, 'd').hours()
      this.minutes = d.subtract(this.days, 'h').minutes()
    }
  },
  methods: {
    submitTime: async function () {

      // return
      if (!this.$refs["formDelayTime"].validate()) {
        return
      } else {
        this.$refs["formDelayTime"].resetValidation()
      }

      let Nanosecond   = 1
      let Microsecond  = 1000 * Nanosecond
      let Millisecond  = 1000 * Microsecond
      let Second       = 1000 * Millisecond
      let Minute       = 60 * Second
      let Hour         = 60 * Minute
      let Day          = 24 * Hour
      let Week         = 7 * Day

      // let res = Second * this.minutes * Minute + this.hours * Hour + this.days * Day + this.weeks * Week
      let res = this.minutes * Minute + this.hours * Hour + this.days * Day + this.weeks * Week
      this.$emit('choice', res)
    }
  },

}
</script>

<style lang="scss">
  .time-delay {
    display: grid;
    grid-template-columns: 1fr 1fr;
    column-gap: 40px;
    & > .time-delay-item {
      display: grid;

      grid-template-columns: 1fr 1fr;

      justify-content: center;
      column-gap: 20px;
      text-align: center;
      & > h3 {
        font-weight: 500;
        font-size: 14px;
        color: cadetblue;
      }

    }

  }
</style>
