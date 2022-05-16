<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">

          <h1 style="font-size: 26px;">{{ isNewOrder ? 'Создания заказа' : 'Редактировать заказ' }}
            <small v-if="item.id">
              <i style="font-size: 14px;">#</i>
              <em style="color: #e02c80;">{{ isNewOrder ? '' : $route.params.public_id }}</em>
              <i style="font-size: 16px;padding-left: 2px;padding-right: 2px;">/</i><em style="color: #0b8484;">{{ item.id }}</em>
            </small>
          </h1>

          <v-btn @click="$router.push({name: 'orders.list'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку заказов</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat>

          <!-- Создание заказа -->
          <v-btn @click="createItem" color="cyan darken-4" class="white--text" small v-if="isNewOrder">
            <v-icon class="mr-3" small>fal fa-save</v-icon>
            Создать заказ
          </v-btn>

          <template v-else>
            <v-btn @click="uploadItem(true)" color="secondary" small>
              <v-icon class="ml-1 mr-3" small>far fa-sync</v-icon>
              Обновить данные
            </v-btn>
            <v-divider class="mx-2 mt-0" inset vertical></v-divider>

            <v-btn @click="updateItemData" color="cyan darken-4" dark small>
              <v-icon class="mr-3" small>fal fa-save</v-icon>
              Сохранить
            </v-btn>
          </template>

        </v-card>
      </v-toolbar>

      <v-form model="valid" class="mx-4 mt-6 mb-4" v-if="item">

        <!-- Основные данные -->
        <v-row class="pt-4">
          <v-col cols="5">

            <v-select
                v-model="item['web_site_id']" :disabled="isCompleted"
                class="mt-0" @change="uploadDeliveries"
                :items="webSites" item-text="hostname" item-value="id"
                label="Укажите магазин" aria-required="true">
            </v-select>

            <v-select v-model="item['order_channel_id']" :items="orderChannels" item-text="name" item-value="id" label="Канал" class="mt-0" :disabled="isCompleted"></v-select>

            <v-select
                v-if="item" v-model="item['individual']"
                label="Тип заказчика" class="mt-0" :disabled="isCompleted"
                :items="[{name:'Частное лицо', status:true},{name:'Юридическое лицо', status:false}]"
                item-text="name" item-value="status">
            </v-select>

          </v-col>
          <v-col cols="5" offset="2">

            <!-- Менеджер -->
            <div class="mt-0">

              <h4 class="mt-0 mb-1" style="font-weight: 500;color:#4b696d;">Менеджер
                <small v-if="item.manager_id" class="ml-2" style="font-size: 14px;color: #ffffff;background-color: #99a09b;padding: 3px 5px 0;border-radius: 3px;">
                  Id: {{ item.manager_id }}
                </small>
              </h4>

              <v-btn text small
                     :disabled="isCompleted"
                     :class="{'blue-grey--text': !item.manager || !item.manager_id}"
                     @click="dialogChoiceManager = true"
                     class="item-editable-popup">
                <v-icon class="mr-2" x-small>
                  far fa-newspaper
                </v-icon>
                {{ (item.manager && item.manager_id) ? item.manager.name + ' ' + item.manager.surname : 'Назначить менеджера' }}
              </v-btn>

            </div>

            <!-- Статус заказа -->
            <v-select
                @change="updateOrderStatus(item, $event)"
                :disabled="item.id === 0"
                v-model="item.status_id"
                item-text="name"
                item-value="id"
                dense class="mb-6 mt-6"
                style="max-width: 260px"
                :items="orderStatuses">
              <template v-slot:selection="{ obj, index }">
                <span
                    :class="$api.orderStatus.getStatusColorByGroup(orderStatuses.find(el=>el.id===item.status_id).group)"
                    style="padding: 2px 6px 1px;border-radius: 2px;"
                >{{ orderStatuses.find(el => el.id === item.status_id).name }}</span>
              </template>
            </v-select>

            <!-- Способ олпаты -->
            <v-select
                v-if="delivery"
                :items.sync="paymentMethods" @change="changePaymentMethod"
                item-text="name" item-value="id" return-object
                :disabled="isCompleted"
                style="max-width: 420px;display: inline-block;" class="mb-0 mt-1 pt-0"
                label="Способ оплаты" :value="item.payment_method_id">
            </v-select>

          </v-col>
        </v-row>

        <v-divider class="mb-4 mt-1"></v-divider>

        <!-- Доставка / Заказчик -->
        <v-expansion-panels class="px-0 mx-0 mt-3" multiple hover>
          <v-expansion-panel>
            <v-expansion-panel-header class="blue-grey lighten-4"><span
                style="font-weight: 500;"
                v-text="isDelivery ? 'Доставка [order]' : 'Доставка [план]'"></span></v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-row>
                <v-col cols="12">
                  <v-card class="d-flex align-center pb-2 justify-space-between mt-2" flat style="max-width: 100%;">
                    <h2></h2>
                    <div>
                      <v-btn
                          v-if="delivery"
                          :disabled="(!delivery || !delivery.id) || isCompleted"
                          small text outlined elevation="1" class=""
                          @click="updateDeliveryCost()">
                        <v-icon class="mr-3 orange--text" small>fal fa-repeat</v-icon>
                        Пересчитать стоимость
                      </v-btn>
                      <v-divider
                          class="mx-4 mt-2"
                          inset
                          vertical
                      ></v-divider>
                      <v-btn
                          v-if="delivery"
                          :disabled="(!delivery || !delivery.id) || isCompleted || (delivery.address_required && !item.delivery_order['address'])"
                          small text outlined elevation="1" class="ml-4"
                          @click="updateDeliveryData()">
                        <v-icon class="mr-3 blue--text" small>fal fa-save</v-icon>
                        {{item.delivery_order.id ? 'Сохранить доставку':'Создать доставку'}}
                      </v-btn>
                    </div>
                  </v-card>
                </v-col>
                <v-col cols="5">

                  <!-- Статус доставки -->
                  <template v-if="item && item.id > 0">
                    <v-select
                        v-if="item && isDelivery"
                        label="Статус доставки"
                        :items="deliveryStatuses" :disabled="isDelivered"
                        @change="updateDeliveryStatus(item.delivery_order, $event)"
                        dense item-text="name" item-value="id"
                        style="max-width: 260px" class="mb-2 mt-4"
                        :value="item.delivery_order.status_id"
                    >
                      <template v-slot:selection="{ obj, index }">
                        <span v-if="item.delivery_order && item.delivery_order.status_id"
                          :class="getDeliveryStatusColorByGroup()" class="mt-2 mb-1"
                          style="padding: 2px 6px 1px;border-radius: 2px;"
                          >
                          {{ item.delivery_order.status.name }}
                        </span>
                        <span v-else class="body-2 grey--text mt-2">Создаем задание на доставку...</span>
                      </template>
                    </v-select>
                  </template>

                  <transition name="fade">
                    <v-text-field
                        messages="Почтовый индекс" dense class="mt-0 pt-0 mb-3"
                        v-if="delivery && delivery['address_required']"
                        v-model="item.delivery_order['postal_code']" disabled></v-text-field>
                  </transition>

                  <v-text-field
                      messages="Стоимость доставки" dense class="mt-0 pt-0 mb-3"
                      v-if="item.delivery_order"
                      v-model.number="item.delivery_order.cost" disabled></v-text-field>

                  <v-text-field
                      messages="Вес (грамм)" dense class="mt-0 pt-0 mb-3"
                      v-if="isDelivery && item.delivery_order"
                      :value="item.delivery_order.weight" disabled></v-text-field>
                </v-col>

                <v-col v-if="true" cols="5" offset="2">

                  <v-select
                      v-if="item.customer_id > 0"
                      :items.sync="deliveries" @change="changeDelivery"
                      item-text="name" item-value="code" return-object
                      style="max-width: 420px;display: inline-block;" class="mb-0 mt-4 pt-0" :disabled="isCompleted"
                      label="Способ доставки" :value="item.delivery_order ? item.delivery_order.code : ''">
                  </v-select>
                  <h4 v-else>Чтобы создать доставку - выберите заказчика</h4>

                  <template v-if="delivery">

                    <transition name="fade">
                      <div class="mt-2" v-if="delivery && (delivery.address_required || delivery.postal_code)">
                        <div class="label-text" style="margin-bottom: 4px;">Адрес доставки (по карте)</div>
                        <vue-dadata
                            token="0ecfda0b057aaf5546432d46e1216b3c4efa9cd4"
                            placeholder="Укажите адрес доставки" :disabled="isCompleted"
                            :class="isValidDeliveryAddress ? 'complete' : 'need-complete'"
                            :query="item.delivery_order ? item.delivery_order['address'] : ''"
                            :onChange="setDeliveryAddress">
                        </vue-dadata>
                      </div>
                    </transition>

                    <!-- comment -->
                    <transition name="fade">
                      <v-textarea
                          label="Комментарий" class="mt-6 pt-0 mb-3" outlined rows="4"
                          v-if="delivery && delivery['code']" :disabled="isCompleted"
                          v-model="item.delivery_order.comment"></v-textarea>
                    </transition>

                  </template>
                </v-col>
              </v-row>
            </v-expansion-panel-content>
          </v-expansion-panel>
          <v-expansion-panel>
            <v-expansion-panel-header class="lime lighten-4">
              <span v-if="item.customer_id" class="" style="font-size: 14px;color:#313131;">
                <span v-if="item.customer && !item.customer.is_unknown" style="font-weight: 500;">Заказчик Id: [{{ item.customer_id }}]</span>
                <span v-else style="font-weight: 500;color:#3c4b52;">Заказчик: [Неизвестный]</span>
              </span>
              <span v-else>Заказчик [Не установлен]</span>
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-row v-if="item.customer">
                <v-col cols="5">
                  <v-card class="d-flex align-center pb-2 justify-space-between mt-0" flat style="max-width: 100%;">
                    <v-switch v-model="createUserForm"
                              :label="createUserForm ? 'Форма создания' : 'Форма редактирования'"></v-switch>
                    <v-btn small text outlined v-if="createUserForm"
                           elevation="1" :disabled="isCompleted"
                           @click="createCustomer(true)">
                      <v-icon class="mr-3 green--text" small>far fa-plus</v-icon>
                      Создать клиента
                    </v-btn>
                  </v-card>
                </v-col>

                <v-col cols="6" offset="1">
                  <div class="d-flex justify-end align-center mt-4">
                    <!-- Change customer -->
                    <v-btn small text outlined
                           elevation="1" :disabled="isCompleted"
                           @click="dialogChoiceCustomer = true">
                      <v-icon class="mr-3 blue-grey--text" small>far fa-search</v-icon>
                      <!--{{ item.customer_id ? 'Изменить заказчика' : 'Найти клиента' }}-->
                      Поиск..
                    </v-btn>
                    <v-divider class="mx-2 mb-1" inset vertical></v-divider>

                    <!-- Set Unknown Customer -->
                    <v-btn small text outlined
                           elevation="1" :disabled="isCompleted || this.isNewOrder"
                           @click="setUnknownCustomer">
                      <v-icon class="mr-3 brown--text text--lighten-2" small>far fa-snowman</v-icon>
                      Неизвестный
                    </v-btn>

                    <v-divider class="mx-2 mb-1" inset vertical></v-divider>
                    <!-- Save data of customer -->
                    <v-btn small text outlined :disabled="!item.customer_id || this.isNewOrder"
                           elevation="1"
                           @click="updateCustomer(true, true)">
                      <v-icon class="mr-3 blue--text" small>far fa-save</v-icon>
                      Сохранить данные
                    </v-btn>
                  </div>
                </v-col>

                <v-col cols="5">

                  <div
                      v-if="item.customer.is_unknown" class="mb-4"
                      style="padding: 4px 16px 3px;border: 1px dotted #009688;display: inline-block;font-weight: 400;">Неизвестный клиент</div>

                  <!-- Редактировать данные -->
                  <template v-if="!createUserForm">
                    <!-- Телефон -->
                    <v-text-field
                        label="Телефон" :disabled="item.customer.is_unknown"
                        v-model.trim="item.customer.phone"
                        class="mt-0"
                    ></v-text-field>

                    <!-- email -->
                    <v-text-field
                        label="Email" :disabled="item.customer.is_unknown"
                        v-model.trim="item.customer.email"
                        class="mt-0"
                    ></v-text-field>

                    <!-- Фамилия -->
                    <v-text-field
                        label="Фамилия" :disabled="item.customer.is_unknown"
                        v-model.trim="item.customer.surname"
                        class="mt-0"
                    ></v-text-field>

                    <!-- Имя -->
                    <v-text-field
                        label="Имя" :disabled="item.customer.is_unknown"
                        v-model.trim="item.customer.name"
                        class="mt-0"
                    ></v-text-field>

                    <!-- Отчество -->
                    <v-text-field
                        label="Отчетство" :disabled="item.customer.is_unknown"
                        v-model.trim="item.customer.patronymic"
                        class="mt-0"
                    ></v-text-field>
                  </template>

                  <template v-else>
                    <!-- Телефон -->
                    <v-text-field
                        label="Телефон"
                        v-model.trim="user.phone"
                        :error-messages="user.errors.phone"
                        class="mt-0"
                    ></v-text-field>

                    <!-- email -->
                    <v-text-field
                        label="Email"
                        v-model="user.email"
                        :error-messages="user.errors.email"
                        class="mt-0"
                    ></v-text-field>

                    <!-- Фамилия -->
                    <v-text-field
                        label="Фамилия"
                        v-model.trim="user.surname"
                        :error-messages="user.errors.surname"
                        class="mt-0"
                    ></v-text-field>

                    <!-- Имя -->
                    <v-text-field
                        label="Имя"
                        v-model.trim="user.name"
                        :error-messages="user.errors.name"
                        dense class="mt-0"
                    ></v-text-field>

                    <!-- Отчество -->
                    <v-text-field
                        label="Отчетство"
                        v-model.trim="user.patronymic"
                        :error-messages="user.errors.patronymic"
                        class="mt-0"
                    ></v-text-field>

                    <!-- Подписка на рассылку -->
                    <v-checkbox v-model="user.subscribed" label="Подписать на рассылку"></v-checkbox>
                  </template>


                </v-col>

                <v-col cols="6" offset="1">
                  <h2 class="text-h3 pb-3">Комментарий от клиента</h2>
                  <v-textarea v-model="item['customer_comment']" class="rounded body-2" outlined placeholder="Оставьте комментарий..." rows="8"></v-textarea>
                </v-col>

              </v-row>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>

        <!-- Cart Items -->
        <v-data-table
            v-if="item && item.id && item.cart_items"
            :headers="cartItemHeaders"
            :items.sync="item.cart_items"
            sort-by="id" item-key="id"
            class="elevation-0 mt-6"
        >
          <template v-slot:top>
            <v-toolbar color="white" flat>

              <v-toolbar-title>
                <h2 style="font-size: 21px;">Список товаров <small>({{item.cart_items.length}})</small></h2>
              </v-toolbar-title>

              <v-divider
                  class="mx-8 mt-2"
                  inset
                  vertical
              ></v-divider>

              <div class="d-flex flex-column">
                <div class="cyan--text text--darken-3 body-1">Скидка: <span style="color: #616161;">{{totalDiscount.toLocaleString('ru-RU',{currency:"RUB"})}} руб.</span></div>

                <div class="green--text text--darken-4 body-1">Итого к оплате:
                  <show-price :price="item.cost" color="grey--text text--darken-2"></show-price>
                </div>

              </div>

              <v-spacer></v-spacer>

              <v-card class="d-flex align-center" flat>
                <v-btn small outlined text @click="setReserveOrder" elevation="1" style="font-weight: 400;color:#101928;" :disabled="isCompleted">
                  <v-icon class="mr-3 mк-2 teal--text text--darken-2" x-small>far fa-cubes</v-icon>Зарезервировать
                </v-btn>
                <v-divider class="mx-4 mb-1" inset vertical></v-divider>
                <v-btn small outlined text @click="dialogChoiceProduct = true" elevation="1" style="font-weight: 400;color:#101928;" :disabled="isCompleted">
                  <v-icon class="mr-3 mк-2 green--text" x-small>far fa-plus</v-icon>Добавить товар
                </v-btn>
              </v-card>

            </v-toolbar>
          </template>

          <template v-slot:item.product.public_id="{ item }">
            {{item.product.public_id}}/{{item.product.id}}
          </template>

          <template v-slot:item.product.label="{ item }">
            <div v-if="item.product_id !== 0" style="max-width: 220px;">
              {{item.product.label}} - <span class=" blue-grey--text">[{{item.product.short_label}}]</span>
            </div>
            <div v-else style="max-width: 220px;">
              {{item.description}}
            </div>
          </template>

          <template v-slot:item.measurement_unit="{ item }">
            <span v-if="item.product_id !== 0">{{item.product['measurement_unit'] ?  (item.product['measurement_unit'].name.toLowerCase().trim()  + ' (' + item.product['measurement_unit'].short_name + ')') : '-'}}</span>
            <span v-else>-</span>
          </template>

          <template v-slot:item.product.article="{ item }">
            <span v-if="item.product_id !== 0">{{item.product.article}}</span>
            <span v-else>-</span>
          </template>

          <template v-slot:item.product.product_tags="{ item }">
            <show-tags :tags="item.product.product_tags"></show-tags>
          </template>

          <!-- Warehouse label -->
          <template v-slot:item.warehouse_item_id="{ item }">
            <div
                v-if="item.product_id > 0"
                @click="openChoiceWarehouseItem(item, availabilityInWarehouses(item))"
                class="item-editable-popup body-2"
            >
              <template v-if="item.reserved || item.wasted" >
                <div :class="item.wasted || item.reserved ? 'grey--text text--darken-2' : 'green--text text--darken-1'">
                  {{getWarehouseLabel(item.warehouse.id)}}
                </div>
              </template>
              <template v-else>
                <div v-if="availabilityInWarehouses(item)" class="lime--text text--darken-3">
                    <div v-if="item.warehouse && item.warehouse.id" >
                      {{getWarehouseLabel(item.warehouse.id)}} - <span class="blue-grey--text">[{{ item.warehouse_items.length }}]</span>
                    </div>
                    <span v-else class="brown--text">на
                    {{item.warehouse_items.length}}
                    {{DeclOfNum(item.warehouse_items.length,['складе', 'складах', 'складах'])}}</span>
                </div>
                <div v-else class="deep-orange--text text--lighten-1">нет в наличии</div>
              </template>
            </div>
            <div v-else>-</div>
          </template>

          <template v-slot:item.reserved="{ item }">
            <div v-if="item.product_id > 0">
              <v-btn x-small icon :disabled="!item.warehouse || !item.warehouse.id || !availabilityInWarehousesReserve(item)" @click="changeReservedCartItem(item)">
                <v-icon :class="item['reserved'] ? 'green--text text--darken-2' : 'deep-orange--text text--darken-1'">far fa-cubes</v-icon>
              </v-btn>
            </div>
            <div v-else>-</div>
          </template>

          <template v-slot:item.wasted="{ item }">
            <div v-if="item.product_id > 0">
              <v-btn  x-small icon :disabled="!item.warehouse || !item.warehouse.id || !availabilityInWarehousesWasted(item)" @click="changeWastedCartItem(item)">
                <v-icon :class="item['wasted'] ? 'green--text text--darken-1' : 'deep-orange--text text--darken-1'">far fa-cubes</v-icon>
              </v-btn>
            </div>
            <div v-else>-</div>
          </template>

          <template v-slot:item.quantity="props" >
            <v-edit-dialog
                v-if="props.item.product_id > 0"
                :return-value.sync="props.item.product.quantity"
                @save="saveAmountUnits(props.item)"
            >
              <div class="text-end item-editable-popup" style="min-width: 80px;" :class="{'grey--text text--darken-1': (isCompleted || props.item.reserved || props.item.wasted)}">
                {{ props.item.product_id === 0 ? props.item.quantity.toFixed(0) : props.item.product.is_integer ?
                  (props.item.quantity ? props.item.quantity.toFixed(0) : '0') : (props.item.quantity ? props.item.quantity.toFixed(2) : '0') }}
                <small>{{props.item.product['measurement_unit'].short_name}}</small>
              </div>

              <template v-slot:input v-if="!props.item.reserved && !props.item.wasted">
                <v-currency-field
                    :disabled="(isCompleted || props.item.reserved || props.item.wasted)"
                    :value-as-integer="props.item.product_id === 0 || props.item.product.is_integer"
                    :decimal-length="props.item.product_id === 0 ? 0 : (props.item.product.is_integer ? 0 : 2)"
                    label="ед. в заказе"
                    v-model.number="props.item.quantity">
                </v-currency-field>
              </template>
            </v-edit-dialog>
            <div v-else>-</div>
          </template>

          <template v-slot:item.product.retail_discount="{ item }">
              <span v-if="item.product_id !== 0">
                {{item.product['retail_discount'] ? item.product['retail_discount'].toLocaleString('ru-RU',{currency:"RUB"}) + 'руб.' : '-'}}
              </span><span v-else>-</span>

          </template>

          <template v-slot:item.product.retail_price="{ item }">
              <span v-if="item.product_id !== 0">
              {{item['product']['retail_price'] ? item['product']['retail_price'].toLocaleString('ru-RU',{currency:"RUB"}) : '-'}} руб.
              </span><span v-else>-</span>
          </template>

          <template v-slot:item.cost="{ item }">
            {{ (item['cost']).toLocaleString('ru-RU',{currency:"RUB"})}} руб.
          </template>

          <template v-slot:item.total_cost="{ item }">
            {{ (item['cost']*item['quantity']).toLocaleString('ru-RU',{currency:"RUB"}) }} руб.
          </template>

          <template v-slot:item.actions="{ item }">
            <show-action :route="{name:'product.edit', params: {public_id:item.product['public_id']}}"
                         :disabledRemove="isCompleted"
                         @delete="removeProduct(item)" removeIcon></show-action>
          </template>

        </v-data-table>
      </v-form>

      <choice-user v-if="item" :open.sync="dialogChoiceManager" @choice="choiceManager" @close="dialogChoiceManager = false"></choice-user>
      <choice-user v-if="item" :open.sync="dialogChoiceCustomer" @choice="choiceCustomer" @close="dialogChoiceCustomer = false"></choice-user>
      <choice-product v-if="item" :open.sync="dialogChoiceProduct" @choice="choiceProduct" @close="dialogChoiceProduct = false"></choice-product>
      <choice-warehouse v-if="item" :open.sync="dialogChoiceWarehouseCartItem" @choice="choiceWarehouseCartItem" 
                        @close="dialogChoiceWarehouseCartItem = false" :wIdx.sync="choicedCartItemWh" :fIdx.sync="reservedWarehouseListIdx"></choice-warehouse>

    </v-card>

    <!-- Подтверждение удаления -->
    <v-dialog
        @keydown.esc="dialogCartItem.open = false"
        max-width="380"
        v-model="dialogCartItem.open"
    >
      <v-card>
        <v-card-title style="font-size: 21px;">Удалить из списка товаров?</v-card-title>

        <v-card-text class="d-flex flex-column">
          <v-list v-if="dialogCartItem && dialogCartItem.product">
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.id">
              <v-list-item-title class="headline">Id: &laquo;{{ dialogCartItem.id }}&raquo;</v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.product.name">
              <v-list-item-title class="headline">Имя: &laquo;{{ dialogCartItem.product.name }}&raquo;
              </v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.product.model">
              <v-list-item-title class="headline">Модель: &laquo;{{ dialogCartItem.product.model }}&raquo;
              </v-list-item-title>
            </v-list-item>
            <v-list-item class="px-0 mx-0 mb-0 pb-0" v-if="dialogCartItem.product.sku">
              <v-list-item-title class="headline">SKU: &laquo;{{ dialogCartItem.product.sku }}&raquo;
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>

          <!-- Отменить -->
          <v-btn
              @click.stop="dialogCartItem.open = false"
              color="green darken-1"
              small
              text
          >
            Отменить
          </v-btn>

          <!-- Удалить -->
          <v-btn
              @click.stop="deleteCartItemDialog"
              color="red darken-1"
              small
              text
          >
            Удалить
          </v-btn>

        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
import draggable from "vuedraggable"

const preloads = 'CartItems.Warehouse,CartItems.Product,CartItems.Product.MeasurementUnit,CartItems.Product.ProductCards,CartItems.Product.ProductTags,CartItems.Product.PaymentSubject,Customer,Manager,OrderChannel,Payment,PaymentMethod,Status,DeliveryOrder,DeliveryOrder.Status';

export default {
  data: () => ({
    item: {
      delivery_order: {
        address: null,
        comment: null,
        postal_code: null,
        code: '',
      },
    },
    user: {
      email: null,
      phone_region: 'IN',
      phone: null,
      name: null,
      surname: null,
      patronymic: null,
      subscribed: true,
      errors: {
        email: null,
        phone: null,
        name: null,
        surname: null,
        patronymic: null,
      },
    },
    orderChannels: [],
    webSites: [],
    deliveryStatuses: [],
    orderStatuses: [],
    deliveries: [],
    warehouses: [],
    dialogChoiceManager: false,
    dialogChoiceCustomer: false,
    dialogChoiceProduct: false,
    dialogChoiceWarehouseCartItem: false,
    choiceCartItem: null, // выбор склада для cartItem
    dialogCartItem: {
      open: false,
    },
    cartItemHeaders: [
      // { text: 'Id', align: 'start', value: 'id'},
      { text: 'p.Id', align: 'start', value: 'product.public_id', sortable: false},
      { text: 'Наименование', value: 'product.label', sortable: false},
      { text: 'Артикул', align: 'end', value: 'product.article', sortable: false},
      // { text: 'Теги', align: 'end', value: 'product.product_tags', sortable: false},
      { text: 'Ед. измерения', align: 'end', value: 'measurement_unit', sortable: false},
      // { text: 'Сборный товар', align: 'end', value: 'product.is_kit'},
      { text: 'Цена', align: 'end', value: 'product.retail_price'},
      { text: 'Скидка', align: 'end',value: 'product.retail_discount'},
      { text: 'Стоимость', align: 'end',value: 'cost'},
      { text: 'Кол-во', align: 'end',value: 'quantity'},
      { text: 'Склад списания', align: 'end',value: 'warehouse_item_id'},
      { text: 'Резерв', align: 'center',value: 'reserved'},
      { text: 'Wasted', align: 'center',value: 'wasted'},
      { text: 'Итого', align: 'end',value: 'total_cost'},


      { text: 'Действие', align: 'center', value: 'actions', sortable: false},
    ],
    createUserForm: false,
  }),
  computed: {
    isNewOrder() {
      return !this.$route.params.public_id
    },
    isDelivery() {
      return !(!this.item.delivery_order || !this.item.delivery_order.id);
    },
    isDelivered() {
      return (this.item.delivery_order && this.item.delivery_order.status_id === 5);
    },
    isCompleted() {
      return (this.item.status_id && this.item.status_id === 10);
    },
    // current delivery
    delivery() {
      if (!this.item.delivery_order || !this.item.delivery_order.code) return null
      return this.deliveries.find((el) => el.code === this.item.delivery_order.code)
    },
    totalCost(){
      let total = 0
      this.item.cart_items.forEach(el=>{
        total += el.cost*el.quantity
      })
      return total
    },
    totalDiscount(){
      let total = 0
      this.item.cart_items.forEach(el=>{
        total += el.product.retail_discount * el.quantity
      })
      return total
    },
    paymentMethods(){
      if (!this.delivery) return []
      return this.delivery.payment_methods
    },
    choicedCartItemWh() {
      if (!this.choiceCartItem || !this.choiceCartItem.warehouse_item) return []
      let arr = []
      arr.push(this.choiceCartItem.warehouse_item.warehouse_id)
      return arr
    },
    reservedWarehouseListIdx() {

      // Собираем Id складов, где вообще можно зарезервировать позицию в нужном объеме
      if (!this.choiceCartItem || !this.choiceCartItem.warehouse_items) return []
      let arr = []
      arr = _.map(this.choiceCartItem.warehouse_items,'warehouse_id')
      return arr
    },
  },
  async mounted() {

    this.createUserForm = !!this.isNewOrder;

    await Promise.all([
      this.uploadOrderChannels(),
      this.uploadOrderStatuses(),
      this.uploadDeliveryStatuses(),
      this.uploadWebSites(),
      this.uploadDeliveries(),
      this.uploadWarehouses(),
    ]).then(() => {

      if (this.isNewOrder) this.item = {
        id: 0,
        accountHashId: this.$store.state.account.hash_id,
        cost: 0,
        manager_id: this.$store.state.user.profile.id,
        manager: this.$store.state.user.profile,
        web_site_id: this.webSites.length > 0 ? this.webSites[0].id : null,
        individual: true,
        order_channel_id: 1,
        customer_comment: "",
        customer: {
          phone: '',
          email: '',
          name: '',
          surname: '',
          patronymic: '',
          issuer_account_id: this.$store.state.account.id,
        },
        payment_method_id: this.paymentMethods.length > 0 ? this.paymentMethods[0].id : null,
        payment_method_type: this.paymentMethods.length > 0 ?this.paymentMethods[0].type : null,
        cart_items: [],
        status_id: 1,
      }
      else this.uploadItem()
    })
  },
  methods: {

    uploadOrderChannels: async function (showNotification = false) {

      this.loading = true
      return this.$api.orderChannel.getListPagination({
        accountHashId: this.$store.state.account.hash_id,
      })
          .then(resp => {
            if (resp['order_channels'] !== undefined) {
              this.orderChannels = resp['order_channels']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring',
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Шаблон не найден',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)
    },
    uploadOrderStatuses: async function (showNotification = false) {

      return this.$api.orderStatus.getList({
        accountHashId: this.$store.state.account.hash_id,
      })
          .then(resp => {
            if (resp['order_statuses'] !== undefined) {
              this.orderStatuses = resp['order_statuses']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring',
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Заказ не найден',
              text: err['message'],
              type: 'error',
            });
          })
    },
    uploadDeliveryStatuses: async function (showNotification = false) {

      return this.$api.deliveryStatuses.getList({
        accountHashId: this.$store.state.account.hash_id,
      })
          .then(resp => {
            if (resp['delivery_statuses'] !== undefined) {
              this.deliveryStatuses = resp['delivery_statuses']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring',
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Заказ не найден',
              text: err['message'],
              type: 'error',
            });
          })
    },
    uploadDeliveryByWebSite: async function (showNotification = false) {

      return this.$api.deliveryStatuses.getList({
        accountHashId: this.$store.state.account.hash_id,
      })
          .then(resp => {
            if (resp['delivery_statuses'] !== undefined) {
              this.deliveryStatuses = resp['delivery_statuses']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring',
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Заказ не найден',
              text: err['message'],
              type: 'error',
            });
          })
    },
    uploadWebSites: async function (showNotification = false) {
      await this.$api.webSite.getList({accountHashId: this.$store.state.account.hash_id})
          .then(resp => {
            if (resp['web_sites'] !== undefined) {
              this.webSites = resp['web_sites']
              if (this.webSites.length > 0) {
                this.item.web_site_id = this.webSites[0].id
                this.uploadDeliveries()
              }
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: 'Ошибка в ответе сервера: web_sites - not found',
                  type: 'warring',
                });
              }
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)
    },
    uploadDeliveries: async function (showNotification = false) {
      if (!this.item.web_site_id || (this.item.web_site_id < 1)) return;
      await this.$api.deliveries.getList({
        accountHashId: this.$store.state.account.hash_id,
        web_site_id: this.item.web_site_id,
      })
          .then(resp => {
            if (resp['deliveries'] !== undefined) {
              this.deliveries = resp['deliveries']
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: 'Ошибка в ответе сервера: deliveries - not found',
                  type: 'warring',
                });
              }
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)
    },
    uploadWarehouses: async function (showNotification = false) {
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
      }

      this.$api.warehouse.getList(payload)
          .then(resp => {
            if (resp['warehouses'] !== undefined ) {
              this.warehouses = resp['warehouses']
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Ошибка в ответе сервера: warehouses - not found',
                type: 'warring'
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error'
            });
          })
    },

    createItem: async function () {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        order_channel_id: this.item.order_channel_id,
        web_site_id: this.item.web_site_id,
      }
      Object.assign(payload,this.item)

      await this.$api.order.create(payload)
          .then((resp) => {
            if (resp['order'] !== undefined) {
              this.item = resp['order']
              return this.$router.push({name: 'order.edit', params: {public_id: resp['order'].public_id}})
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [order]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
    },
    createCustomer: async function (showNotification = false) {

      if(!this.item) return

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        // issuer_account_id: this.$store.state.account.id,

        phone: this.user.phone,
        phone_region: this.user.phone_region,
        email: this.user.email,
        name: this.user.name,
        surname: this.user.surname,
        patronymic: this.user.patronymic,

        subscribed: this.user.subscribed,
        role_id: 7, //
      }

      await this.$api.user.create(payload)
          .then(async (resp) => {
            if (resp['user'] !== undefined) {


              // Сохраняем данные нового клиента
              this.item.customer = resp['user']
              this.item.customer_id = resp['user'].id

              if (!this.isNewOrder) await this.updateItemData(false)
              else  this.createUserForm = false

              if (showNotification)  {
                this.$notify({
                  group: 'main',
                  title: 'Клиент создан',
                  type: 'success',
                });
              }

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [user]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.user.errors = err.errors
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {})

    },
    choiceCustomer: async function (user) {
      this.item.customer_id = user.id
      this.item.customer = user

      if (!this.isNewOrder) {
        await this.updateItemData(true)
            .finally(()=>{this.dialogChoiceCustomer = false; this.createUserForm = false;})
      } else {
        this.dialogChoiceCustomer = false;
        this.createUserForm = false;
      }

    },
    updateCustomer: async function (showNotification = false, updateItem = false) {

      if(!this.item || !this.item.customer_id || !this.item.customer) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.customer_id,
        phone: this.item.customer.phone,
        email: this.item.customer.email,
        name: this.item.customer.name,
        surname: this.item.customer.surname,
        patronymic: this.item.customer.patronymic,
        phone_region: 'IN',
      }

      await this.$api.user.update(payload)
          .then((resp) => {
            if (resp['user'] !== undefined) {
              if (showNotification)  {
                this.$notify({
                  group: 'main',
                  title: 'Данные клиента обновлены',
                  type: 'success',
                });
              }


            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [user]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(()=>this.dialogChoiceCustomer = false)

      if (updateItem) await this.uploadItem(false)
    },
    setUnknownCustomer: async function () {

      // Нет смысла создавать, установим неизвестного клиента
      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id, // client
        preloads
      }
      this.$api.order.setUnknownCustomer(payload)
          .then((resp) => {
            if (resp['order'] !== undefined) {

              this.$set(this, 'item', resp['order'])

              this.$notify({
                group: 'main',
                title: 'Заказчик обновлен',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [order]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
    },

    choiceManager: async function (user) {
      this.item.manager_id = user.id
      this.item.manager = user
      this.dialogChoiceManager = false

      // Не обновляем данные, если это новый заказ
      if (!this.item || this.item.id === 0) return

      /* Обновляем только менеджера, другие данные не трогаем */
      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        public_id: this.item.public_id,
        manager_id: user.id,
        preloads,
      }

      await this.$api.order.update(payload)
          .then((resp) => {
            if (resp['order'] !== undefined) {
              this.$notify({
                group: 'main',
                title: 'Менеджер обновлен',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [order]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
    },
    choiceProduct: async function (product) {

      if (!product || !product.id) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_id: product.id,
        preloads
      }

      this.$api.order.appendProduct(payload)
          .then(async (resp) => {
            if (resp['order']) {

              if (!resp['order']['cart_items']) this.item['cart_items'] = []
              else this.item['cart_items'] = resp['order']['cart_items']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе [order]',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка добавления товара',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.dialogChoiceProduct = false)

    },
    removeProduct: async function (cartItem) {

      if (cartItem.product_id === 0) {

        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: cartItem.id,
          preloads
        }

        this.$api.cartItem.delete(payload)
            .then(async () => {
              // удаляем доставку
              let payload = {
                accountHashId: this.$store.state.account.hash_id,
                id: this.item.delivery_order.id,
              }
              await this.$api.deliveryOrder.delete(payload)
                  .then( async () => await this.uploadItem())
                  .catch((err) => {
                    this.$notify({
                      group: 'main',
                      title: 'Ошибка обновления',
                      text: err['message'],
                      type: 'error',
                    });
                  })
            })
            .catch((err) => {
              this.$notify({
                group: 'main',
                title: 'Ошибка создания',
                text: err['message'],
                type: 'error',
              });
            })

      } else {
        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.item.id,
          product_id: cartItem.product_id,
          preloads
        }

        this.$api.order.removeProduct(payload)
            .then(async (resp) => {
              if (resp['order']) {

                if (!resp['order']['cart_items']) this.item['cart_items'] = []
                else this.item['cart_items'] = resp['order']['cart_items']

                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              } else {
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: 'Отсутствует переменная в ответе order',
                  type: 'warring',
                });
              }

            })
            .catch((err) => {
              this.$notify({
                group: 'main',
                title: 'Ошибка создания',
                text: err['message'],
                type: 'error',
              });
            })
            .finally(async ()=> await this.uploadItem())
      }


    },
    saveAmountUnits: async function(item){

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: item.id,
        quantity: item.quantity,
        preloads: "Product"
      }

      this.$api.cartItem.update(payload)
          .then(async (resp) => {
            if (resp['cart_item']) {
              item = resp['cart_item']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе [cart_item]',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка добавления товара',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.uploadItem())
    },

    // Изменение способа доставки
    changeDelivery: async function(o) {
      if (!o) return
      // console.log(e)
      if (!this.item.delivery_order || !this.item.delivery_order.id) {
        this.item.delivery_order = {
          address: '', postal_code: '', total_cost: 0, comment: '', code: o.code, web_site_id: this.item.web_site_id
        }
      }

      this.item.delivery_order.code = o.code

      await this.updateDeliveryCost()
    },
    setDeliveryAddress: async function (data) {
      if (data['value'] !== undefined && data['data']['postal_code'] !== undefined) {
        if (!this.item.delivery_order) this.item.delivery_order = {address: '', postal_code: '', total_cost: 0, comment: ''}
        this.item.delivery_order.address = data.value
        this.item.delivery_order.postal_code = data.data.postal_code
        await this.updateDeliveryCost()
      } else {
        this.item.delivery_order.postal_code = null
        // this.calculatedDelivery = false // сбрасываем успешную стоимость расчета
        // this.deliveryCost = 0 // устанавливаем нулевую стоимость доставки
        this.item.delivery_order.total_cost = null // устанавливаем нулевую стоимость доставки
      }

    },
    updateDeliveryCost: async function () {

      if (!this.item.web_site_id || this.item.web_site_id < 1) {
        return
      }
      if (!this.item.delivery_order) this.item.delivery_order = {address: '', postal_code: '', total_cost: 0, comment: ''}

      if (this.delivery['address_required'] && !this.item.delivery_order.address) return
      if (this.delivery['postal_code_required'] && !this.item.delivery_order.postal_code) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        web_site_id: this.item.web_site_id,
      }
      Object.assign(payload, this.getPayloadCart())

      // console.log(this.getPayloadCart().cart)
      // return

      await this.$api.deliveries.calculate(payload)
          .then(resp => {
            if (resp['total_cost'] !== undefined) {
              this.$set(this.item.delivery_order,'cost',  resp['total_cost'])
              this.$set(this.item.delivery_order,'weight',resp['weight'])
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Ошибка в ответе сервера: total_cost - not found',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)


    },
    getPayloadCart() {
      if (!this.delivery) return
      let deliveryData = {}
      deliveryData.id = this.delivery.id
      deliveryData.code = this.delivery.code
      deliveryData.address_required = this.delivery.address_required
      deliveryData.postal_code = this.item.delivery_order.postal_code

      if (this.delivery.payment_method) {
        deliveryData.payment_method = this.delivery.payment_method
      }

      let cart = this.item.cart_items;

      // новая корзина для убирания всех не нужных данных
      let cartItems = [];

      // снижаем объем передаваемых данных
      for (let key in cart) {
        // console.log(key)
        cartItems.push({
          id: cart[key].product_id,
          quantity: cart[key].quantity,
        })
      }

      return {
        cart: cartItems,
        delivery_data: deliveryData,
      };
    },
    isValidDeliveryAddress() {
      if (!this.delivery || this.delivery['address_required'] === undefined) {
        console.log('не выбран метода доставки')
        return false; // не выбран метода доставки
      }
      if (this.delivery['postal_code_required'] === true && (this.delivery.postal_code === undefined || this.delivery.postal_code === null)) {
        console.log('требуется индекс, но его нет')
        return false; // требуется индекс, но его нет
      }
      if (this.delivery['address_required'] === true && (this.delivery.address === undefined || this.delivery.address === null)) {
        console.log('требуется адрес, но его нет')
        return false; //
      }

      return true
    },
    changePaymentMethod: async function(e) {

      this.item.payment_method_id = e.id
      this.item.payment_method_type = e.type

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        payment_method_id: this.item.payment_method_id,
        payment_method_type: this.item.payment_method_type,
        preloads,
      }

      await this.$api.order.update(payload)
          .then((resp) => {
            if (resp['order'] !== undefined) {
              this.$notify({
                group: 'main',
                title: 'Способ оплаты обновлен',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [order]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
    },

    uploadItem: async function (showNotification = false) {

      this.loading = true
      await this.$api.order.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads,
      })
          .then(resp => {
            if (resp['order'] !== undefined) {
              this.item = resp['order']
              if (this.item.delivery_order) {
                let status_id = this.item.delivery_order.status_id
                this.item.delivery_order.status = this.deliveryStatuses.find(el => el.id === status_id, status_id)
              }

              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера получены не полностью',
                type: 'warring',
              });
            }
          })
          .catch(err => {
            this.$notify({
              group: 'main',
              title: 'Заказ не найден',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)

    },
    updateItemData: async function (showNotification = false) {

      if (this.delivery) await this.updateDeliveryData()
      if (this.item.customer && this.item.customer_id) await this.updateCustomer(false)

      let payload = this.item

      Object.assign(payload, {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        preloads,
      })

      await this.$api.order.update(payload)
          .then((resp) => {
            if (resp['order'] !== undefined) {
              this.item = resp['order']
              if (this.item.delivery_order && this.delivery) {
                let status_id = this.item.delivery_order.status_id
                this.item.delivery_order.status = this.deliveryStatuses.find(el => el.id === status_id, status_id)
              }
              if (showNotification) {
                this.$notify({
                  group: 'main',
                  title: 'Даныне заказа обновлены',
                  type: 'success',
                });
              }
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [order]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(()=> this.createUserForm = false)

    },
    openDeleteCartItemDialog: async function (item) {
      Object.assign(this.dialogCartItem, item)
      this.dialogCartItem.open = true
    },
    deleteCartItemDialog: async function () {
      this.dialogCartItem.open = false
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.dialogCartItem.id,
      }

      this.$api.order.update(payload)
          .then(async () => {
            await this.uploadItem(false)
            Object.assign(this.dialogCartItem, {open: false})
            this.$notify({
              group: 'main',
              title: 'Товар удален из списка',
              type: 'success',
            });


          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка удаления',
              text: err['message'],
              type: 'error',
            });
          })


    },

    updateOrderStatus: async function (item, status_id) {
      let id = item.status_id
      let _itemStatus = this.orderStatuses.find(el => el.id === id, id)
      if (_itemStatus) item.status = _itemStatus

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: item.id,
        status_id: status_id,
        preloads,
      }

      await this.$api.order.update(payload)
          .then((resp) => {
            if (resp['order'] !== undefined) {
              this.item = resp['order']
              this.$notify({
                group: 'main',
                title: 'Статус обновлен',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [order]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.uploadItem()
          })
    },
    getDeliveryStatusColorByGroup() {

      let el = this.deliveryStatuses.find(el => el.id === this.item.delivery_order.status_id)

      switch (el.group) {
          // новая доставка
        case "new":
          return "blue  lighten-2 white--text";
          // доставка подтверждена
        case "agreement":
          return "purple  lighten-3 white--text";
          // В процессе
        case "delivery":
          return "brown lighten-2 white--text";
          // Отменена
        case "canceled":
          return "red lighten-3 white--text";
          // Завершена успешно
        case "completed":
          return "green  lighten-1 white--text";
      }
    },
    updateDeliveryData: async function () {

      let newItem = false
      if (!this.item.delivery_order || !this.item.delivery_order.id) newItem = true
      // Если нет существующей доставки и создаем ордер
      if (newItem) {

        // console.log('postal code: ', this.item.delivery_order.postal_code)
        // return
        let payload = {}
        Object.assign(payload, this.item.delivery_order)
        // address: '', postal_code: '', total_cost: 0, comment: '', code: o.code, web_site_id: this.item.web_site_id

        payload = {
          accountHashId: this.$store.state.account.hash_id,
          status_id: 1, // new
          order_id: this.item.id,
          postal_code: this.item.delivery_order ? this.item.delivery_order.postal_code : null, // fix
          customer_id: this.item.customer_id,
          web_site_id: this.item.web_site_id,
          code: this.delivery ? this.delivery.code : null,
          method_id: this.delivery ? this.delivery.id : null,
          cost: this.item.delivery_order ? this.item.delivery_order.cost : null,
          preloads: "Status",
        }

        await this.$api.deliveryOrder.create(payload)
            .then(async  (resp) => {
              if (resp['delivery_order'] !== undefined) {

                this.$set(this.item, 'delivery_order', resp['delivery_order'])

                this.$notify({
                  group: 'main',
                  title: 'Статус доставки обновлен',
                  type: 'success',
                });
              } else {
                this.$notify({
                  group: 'main',
                  title: 'Данные от сервера не полные',
                  text: 'не хватает объекта [delivery_order]',
                  type: 'warring',
                });
              }
              await this.uploadItem()
            })
            .catch((err) => {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: err['message'],
                type: 'error',
              });
            })

      } else {

        let payload = {
          accountHashId: this.$store.state.account.hash_id,
          id: this.item.delivery_order.id,
          preloads: "Status",
        }

        Object.assign(payload,this.item.delivery_order)

        await this.$api.deliveryOrder.update(payload)
            .then(async(resp) => {
              if (resp['delivery_order'] !== undefined) {
                this.item.delivery_order = resp['delivery_order']

                this.$notify({
                  group: 'main',
                  title: 'Данные доставки обновлены',
                  type: 'success',
                });
                await this.uploadItem()
              } else {
                this.$notify({
                  group: 'main',
                  title: 'Данные от сервера не полные',
                  text: 'не хватает объекта [delivery_order]',
                  type: 'warring',
                });
              }
            })
            .catch((err) => {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: err['message'],
                type: 'error',
              });
            })

      }

    },
    updateDeliveryStatus: async function (item, status_id) {

      let id = item.status_id
      let _itemStatus = this.deliveryStatuses.find(el => el.id === id, id)
      if (_itemStatus) item.status = _itemStatus

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: item.id,
        status_id: status_id,
        preloads: "Status",
      }

      await this.$api.deliveryOrder.update(payload)
          .then((resp) => {
            if (resp['delivery_order'] !== undefined) {
              this.item.delivery_order = resp['delivery_order']


              this.$notify({
                group: 'main',
                title: 'Статус доставки обновлен',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [delivery_order]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })

    },

    // Выбор скалада у cart Item
    getWarehouseLabel(warehouse_id) {
      let el = this.warehouses.find(el=>el.id===warehouse_id)
      if (!el) return 'unknown name'

      if (el['short_label'] && el.short_label.length > 1) return el.short_label
      if (el.label && el.label.length > 1) return el.label

      return 'unknown name'

    },
    openChoiceWarehouseItem(item, availability){
      if (!availability || item.reserved || item.wasted) return
      this.dialogChoiceWarehouseCartItem = true
      this.choiceCartItem = item
    },
    availabilityInWarehouses(cartItem){
      if (cartItem.product_id > 0) {
        return cartItem.warehouse_items && cartItem.warehouse_items.length > 0
      }

    },
    availabilityInWarehousesReserve(cartItem){
      if (cartItem.product_id < 1) return true
      if (cartItem.reserved) return true
      return cartItem.warehouse_items && cartItem.warehouse_items.length > 0
    },
    availabilityInWarehousesWasted(cartItem){
      if (cartItem.product_id < 1 ) return true
      if (cartItem.wasted) return true
      if (cartItem.reserved) {
        return true
      } else {
        return cartItem.warehouse_items && cartItem.warehouse_items.length > 0
      }

    },
    choiceWarehouseCartItem: async function (warehouse) {

      if (!warehouse || !warehouse.id || !this.choiceCartItem) return

      // Пытаться изменить склад списания/резерва по warehouse_id
      // Изменять объем заказа, а вместе с ним изменять объем резерва на складе (+ проверка доступного объема)
      // Снимать резерв

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.choiceCartItem.id,
        warehouse_id: warehouse.id, // new warehouse id
        quantity: this.choiceCartItem.quantity,
        reserved: true,
      }

      // console.log(payload)
      // return

      this.$api.cartItem.updateReserve(payload)
          .then(async (resp) => {
            if (resp['cart_item']) {

              this.$notify({
                group: 'main',
                title: 'Данные на складе обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе [cart_item]',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка добавления товара',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.uploadItem())

      this.dialogChoiceWarehouseCartItem = false

    },
    changeReservedCartItem: async function(cartItem) {

      if (!cartItem || !cartItem.id) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: cartItem.id,
        reserved: !cartItem.reserved,
        quantity: cartItem.reserved ? null : cartItem.quantity, // если ставим резерв
        preloads: "WarehouseItem,Product,Product.MeasurementUnit"
      }

      // console.log(payload)
      // return
      this.$api.cartItem.updateReserve(payload)
          .then(async (resp) => {
            if (resp['cart_item']) {
              /*let item = resp['cart_item']
              let index = this.item.cart_items.findIndex(el=> el.id === item.id)
              if (index !== -1) {
                this.$set(this.item.cart_items,0,item)
              }*/
              await this.uploadItem()


              this.$notify({
                group: 'main',
                title: 'Данные на складе обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе [cart_item]',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка добавления товара',
              text: err['message'],
              type: 'error',
            });
          })
          // .finally(() => this.uploadItem())
    },
    changeWastedCartItem: async function(cartItem) {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: cartItem.id,
        wasted: !cartItem.wasted,
        // quantity: cartItem.reserved ? null : cartItem.quantity, // если ставим резерв
        // preloads: "WarehouseItem,Product,Product.MeasurementUnit",
        preloads: "Warehouse,Product,Product.MeasurementUnit"
      }

      // console.log(payload)
      // return
      this.$api.cartItem.updateReserve(payload)
          .then(async (resp) => {
            if (resp['cart_item']) {
              let item = resp['cart_item']
              await this.uploadItem()

              /*let index = this.item.cart_items.findIndex(el=> el.id === item.id, item)
              if (index !== -1) {
                this.$set(this.item.cart_items,index,item)*/

              this.$notify({
                group: 'main',
                title: 'Данные на складе обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе [cart_item]',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка добавления товара',
              text: err['message'],
              type: 'error',
            });
          })
      // .finally(() => this.uploadItem())
    },

    setReserveOrder: async function() {

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        reserved: true,
        preloads,
      }

      await this.$api.order.updateReserve(payload)
          .then((resp) => {
            if (resp['order'] !== undefined) {
              this.item = resp['order']
              this.$notify({
                group: 'main',
                title: 'Статус обновлен',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [order]',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.uploadItem()
          })

    }
  },
  components: {
    draggable,
    ChoiceUser: () => import('@/components/views/tpl/organizations/users/ChoiceUser'),
    ChoiceProduct: () => import('@/components/views/tpl/shops/products/ChoiceProduct'),
    ChoiceWarehouse: () => import('@/components/views/tpl/warehouse/warehouse/ChoiceWarehouse'),
    WorkStatus: () => import('@/components/views/tpl/layouts/WorkStatus'),
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ShowCount: () => import('@/components/views/tpl/layouts/table/ShowCount'),
    ShowTags: () => import('@/components/views/tpl/layouts/table/ShowTags'),
  },
}
</script>
