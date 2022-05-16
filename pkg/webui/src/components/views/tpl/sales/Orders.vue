<template>
    <v-container fluid>
        <h1 class="mb-6 red--text">Заказы (Orders)</h1>

        <v-data-table
                :headers="headers"
                :items="orders"
                class="elevation-1"
                :single-expand="false"
                :expanded.sync="expanded"
                item-key="id"
                :show-expand="false"
                single-expand
                @click:row="clickRow"
                @item-expanded="expandedRow"
        >
            <template v-slot:header.name="{ header }">
                <span class="red&#45;&#45;text">{{ header.text }}</span>
            </template>

            <template v-slot:item.status="{ item }">
                <span class="grey--text">{{ item.status === 'lead' ? 'новый заказ' : 'обрабатывается' }}</span>
            </template>

            <template v-slot:item.cost="{ item }">
                <!--<span :class="{'green&#45;&#45;text':item.payment}">{{ item.payment ? 'оплачен' : 'ожидается' }}</span>-->
                <span class="red--text text--darken-2">{{ item.cost }}</span>
            </template>

            <template v-slot:item.payment="{ item }">
                <!--<span :class="{'green&#45;&#45;text':item.payment}">{{ item.payment ? 'оплачен' : 'ожидается' }}</span>-->
                <span :class="item.payment ? 'green--text' : 'blue--text'">{{ item.payment ? 'оплачен' : 'ожидается' }}</span>
            </template>

            <template v-slot:expanded-item="{ headers, item }" style="box-shadow: none;">
                <td :colspan="headers.length">
                    <div>
                        Some info about order there....
                        More info about {{ item.username }}
                    </div>
                </td>
            </template>

        </v-data-table>
    </v-container>
</template>

<script>

    export default {
        data: () => ({
            drawer: true,
            expanded: [],
            headers: [
                { text: '#', align: 'start', value: 'id' },
                { text: 'Дата и время', value: 'datetime'},
                { text: 'Клиент', value: 'username'},
                { text: 'Адрес доставки', value: 'address'},
                { text: 'Статус', value: 'status'},
                { text: 'Сумма', value: 'cost'},
                { text: 'Оплата', value: 'payment'},
            ],
            orders: [
                {
                    id: 1,
                    datetime: '25/04/2020',
                    username: 'Николай Иванович',
                    address: 'Москва. Коломенский проезд д21, кв253, 10-й подъезд, домофон в253в3689, этаж 2',
                    status: 'lead',
                    cost: 2050,
                    payment: true,

                },
                {
                    id: 2,
                    datetime: '25/04/2020',
                    username: 'Алексей',
                    address: 'г. Ростова-на-Дону, ул. Глинки, д. 2',
                    status: 'lead',
                    cost: 8700,
                    payment: false,
                },
                {
                    id: 3,
                    datetime: '22/04/2020',
                    username: 'Юрий Беликов',
                    address: 'Московская область г Волоколамск СДЕК',
                    status: 'lead',
                    cost: 125.33,
                    payment: true,

                },
                {
                    id: 4,
                    datetime: '23/04/2020',
                    username: 'Зотов Андрей',
                    address: 'проспект Андропова 18к6',
                    status: 'lead',
                    cost: 5800,
                    payment: true,
                },

            ],
        }),
        methods: {
            showUp: function(row){
                console.log(row)
                console.log('Row id: ', row.id)
            },
            clickRow: function(value){

                if (this.expanded.length > 0 && this.expanded[0].id === value.id ) {
                    this.expanded=[]
                } else {
                    this.expanded=[value]
                }
            },
            expandedRow: function(e) {
                let item = e.item
                let value = e.value
              console.log('Expand: ', item)
              console.log('Status: ', value)
            },
        },
        computed: {

        },

    }
</script>

<style lang="scss">
    thead.v-data-table-header tr th span {
        font-size: 14px !important;
        padding-right: 10px;

    }

    .v-data-table__expanded__content {

        box-shadow: unset !important;
        /*background-color: #eeeeee;*/
        background-color: #e8fdee;
    }
    .v-data-table tr:hover {
        cursor: pointer;
    }
</style>