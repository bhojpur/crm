<template>
  <v-container fluid>
    <v-card class="pa-2" style="border-top-left-radius: 2px;border-top-right-radius: 2px;">
      <v-toolbar color="white" flat>
        <v-toolbar-title class="d-flex flex-column pt-6">
          <h1 style="font-size: 26px;">Редактировать товар
            <v-chip v-if="item" color="grey lighten-3">№ {{ item.id }}/{{ item.public_id }}</v-chip>
          </h1>
          <v-btn @click="$router.push({name: 'products.list'})" class="ml-0 pl-0" outlined small text>
            <v-icon class="mr-3 ml-2" small>fal fa-arrow-circle-left</v-icon>
            <span style="padding-top: 2px">Назад к списку</span>
          </v-btn>
        </v-toolbar-title>
        <v-spacer></v-spacer>
        <v-card class="d-flex align-center" flat v-if="item">

          <!-- Обновление и сохранение серии -->
          <v-btn @click="uploadItem(true)" color="secondary" small>
            <v-icon class="ml-1 mr-3" small>far fa-sync</v-icon>
            Обновить данные
          </v-btn>
          <v-divider class="mx-2 mt-0" inset vertical></v-divider>

          <v-btn @click="updateItemData" color="cyan darken-4" class="white--text" small>
            <v-icon class="mr-3" small>fal fa-save</v-icon>
            Сохранить
          </v-btn>

        </v-card>

      </v-toolbar>
      <v-form class="mx-4 mt-6 mb-4" ref="form" v-if="item" v-model="formValid">
        <!-- Основные данные -->
        <v-row class="pt-6">

          <v-col cols="4">

            <!-- Label -->
            <v-text-field
                class="mt-0"
                label="Наименование"
                counter="128"
                v-model.trim="item['label']"
            ></v-text-field>

            <!-- second_label -->
            <v-text-field
                class="mt-0"
                label="Дополнительное наименование"
                counter="128"
                v-model.trim="item['second_label']"
            ></v-text-field>

            <!-- short_label -->
            <v-text-field
                class="mt-0"
                label="Краткое наименование"
                counter="128"
                v-model.trim="item['short_label']"
            ></v-text-field>

            <!-- Артикул -->
            <v-text-field
                class="mt-0"
                label="Артикул"
                counter="128"
                v-model.trim="item['article']"
            ></v-text-field>

            <!-- Модель -->
            <v-text-field
                class="mt-0"
                label="Модель"
                counter="255"
                v-model.trim="item['model']"
            ></v-text-field>

            <!-- Товарный знак -->
            <v-text-field
                class="mt-0"
                label="Товарный знак"
                counter="128"
                v-model.trim="item['trademark']"
            ></v-text-field>

            <!-- brand -->
            <v-text-field
                class="mt-0"
                label="Бренд"
                counter="128"
                v-model.trim="item['brand']"
            ></v-text-field>

            <v-textarea
                label="Краткое описание товара"
                counter="255" rows="2" class="mt-0"
                v-model.trim="item['short_description']"
            ></v-textarea>
          </v-col>

          <v-col cols="4" offset="0" class="pl-8">

            <!-- Единица измерения товара: штуки, метры, литры, граммы и т.д. -->
            <v-select :items="measurementUnits" class="mt-1" item-text="name"
                      item-value="id" label="Единица измерения" v-model="item['measurement_unit_id']"></v-select>

            <!-- Признак предмета расчета (бухучет - № 54-ФЗ) товар, услуга, работа, набор (комплект) = сборный товар -->
            <v-select :items="paymentSubjects" class="mt-1" item-text="name"
                      item-value="id" label="Признак предмета расчета" v-model="item['payment_subject_id']"></v-select>

            <!-- Коды ставок НДС (бухучет) -->
            <v-select :items="vatCodes" class="mt-1" item-text="name"
                      item-value="id" label="Код ставки НДС" v-model="item['vat_code_id']"></v-select>

            <h4 class="mt-1 mb-1 primary--text" style="font-weight: 500;">Категории товара</h4>
            <div @click="dialogChoiceCategories = true" class="d-flex align-center item-editable-popup">
              <v-chip-group column>
                <template v-if="item.product_categories && item.product_categories.length > 0">
                  <v-chip v-for="category in item.product_categories" :key="category.id" class="body-2 blue lighten-5" small>
                    {{ category.label }}
                  </v-chip>
                </template>
                <span v-else class="deep-orange--text">
                  <v-icon class="mr-3 ml-1" small>
                    far fa-folder-tree
                  </v-icon>
                  n/a
                </span>

              </v-chip-group>
            </div>

            <h4 class="mt-4 mb-1 primary--text" style="font-weight: 500;"># Теги товара</h4>
            <div @click="dialogChoiceTags = true" class="d-flex align-center item-editable-popup">
              <v-chip-group column>
                <template v-if="item.product_tags && item.product_tags.length > 0">
                  <v-chip v-for="tag in item.product_tags" :key="tag.id" class="body-2 blue lighten-5" small>
                    {{ tag.label }}
                  </v-chip>
                </template>
                <span v-else class="deep-orange--text">
                  <v-icon class="mr-3 ml-1" small>
                    far fa-tags
                  </v-icon>
                  n/a
                </span>

              </v-chip-group>
            </div>

            <v-divider class="mb-2 mt-4"></v-divider>

            <!-- Тип товара -->
            <div class="d-flex d-inline align-center">
              <h4 class="mt-4 mb-1 primary--text mr-3" style="font-weight: 500;">Тип товара: </h4>
              <v-btn text small
                     :class="{'red--text': !item.product_type_id}"
                     @click="dialogChoiceType = true"
                     class="item-editable-popup">
                <v-icon class="mr-2" x-small>
                  far fa-newspaper
                </v-icon>
                {{ (item.product_type && item.product_type_id) ? item.product_type.name : 'n/a' }}
              </v-btn>
              <v-icon
                  @click="deleteType()"
                  class="blue-grey--text text--lighten-2 ml-3"
                  small
              >
                fas fa-trash
              </v-icon>
            </div>

            <!-- Производитель -->
            <div class="d-flex d-inline align-center mt-2">
              <h4 class="mt-4 mb-1 primary--text mr-3" style="font-weight: 500;">Производитель: </h4>
              <v-btn text small
                     :class="{'deep-orange--text': !item.manufacturer_id}"
                     @click="dialogChoiceManufacturer = true"
                     class="item-editable-popup">
                <v-icon class="mr-2" x-small>
                  far fa-newspaper
                </v-icon>
                {{ getManufacturerName(item) }}
              </v-btn>
              <v-icon
                  @click="deleteManufacturer()"
                  class="blue-grey--text text--lighten-2 ml-3"
                  small
              >
                fas fa-trash
              </v-icon>
            </div>

          </v-col>

          <!-- Цены (розница / опт) -->
          <v-col cols="4" offset="0" class="pl-8">

            <!-- retail_price -->
            <v-currency-field
                label="Розничная цена" class="mt-1"
                :value-as-integer="false"
                v-model.number="item['retail_price']">
            </v-currency-field>

            <!-- retail_discount -->
            <v-currency-field
                label="Розничная скидка"
                :value-as-integer="false"
                v-model.number="item['retail_discount']">
            </v-currency-field>

            <!-- Доступен ли в розницу -->
            <v-checkbox
                v-model="item['enable_retail_sale']"
                class="mb-0 pb-0 my-0" dense
                label="Доступен для розничных продаж"
            ></v-checkbox>

            <!--  Изображения / Атрибуты / Описание  -->
            <v-expansion-panels class="px-0 mx-0 mt-3" multiple hover  >
              <!--<v-expansion-panel>
                <v-expansion-panel-header class="light-green lighten-5">Розничные цены</v-expansion-panel-header>
                <v-expansion-panel-content>
                  &lt;!&ndash; retail_price &ndash;&gt;
                  <v-currency-field
                      label="Розничная цена" class="mt-4"
                      :value-as-integer="false"
                      v-model.number="item['retail_price']">
                  </v-currency-field>

                  &lt;!&ndash; retail_discount &ndash;&gt;
                  <v-currency-field
                      label="Розничная скидка"
                      :value-as-integer="false"
                      v-model.number="item['retail_discount']">
                  </v-currency-field>

                  &lt;!&ndash; Доступен ли в розницу &ndash;&gt;
                  <v-checkbox
                      v-model="item['enable_retail_sale']"
                      class="mb-0 pb-0" dense
                      label="Доступен для розничных продаж"
                  ></v-checkbox>
                </v-expansion-panel-content>
              </v-expansion-panel>-->
              <v-expansion-panel>
                <v-expansion-panel-header class="light-green lighten-5">Оптовые цены</v-expansion-panel-header>
                <v-expansion-panel-content>
                  <!-- wholesale_price_1 -->
                  <v-currency-field
                      label="Оптовая цена - 1" class="mt-4"
                      :value-as-integer="false"
                      v-model.number="item['wholesale_price_1']">
                  </v-currency-field>

                  <!-- wholesale_price_2 -->
                  <v-currency-field
                      label="Оптовая цена - 2"
                      :value-as-integer="false"
                      v-model.number="item['wholesale_price_2']">
                  </v-currency-field>

                  <!-- wholesale_price_3 -->
                  <v-currency-field
                      label="Оптовая цена - 3"
                      :value-as-integer="false"
                      v-model.number="item['wholesale_price_3']">
                  </v-currency-field>

                  <!-- Доступен для оптовых продаж -->
                  <v-checkbox
                      v-model="item['wholesale_sale']"
                      class="my-0 py-0" dense
                      label="Доступен для оптовых продаж"
                  ></v-checkbox>
                </v-expansion-panel-content>
              </v-expansion-panel>
              <v-expansion-panel>
                <v-expansion-panel-header class="light-green lighten-5">Условия хранения</v-expansion-panel-header>
                <v-expansion-panel-content>
                  <!-- manufacture_date -->
                  <!--<v-text-field
                      class="mt-4"
                      label="Дата изготовления"
                      required
                      v-model.trim="item['manufacture_date']"
                  ></v-text-field>-->


                  <!-- storage_requirements -->
                  <v-text-field
                      class="mt-0"
                      label="Условия хранения"
                      required
                      v-model.trim="item['storage_requirements']"
                  ></v-text-field>

                  <!-- shelf_life -->
                  <v-text-field
                      class="mt-0"
                      label="Срок годности"
                      counter="128"
                      v-model.trim="item['shelf_life']"
                  ></v-text-field>
                </v-expansion-panel-content>
              </v-expansion-panel>
              <v-expansion-panel>
                <v-expansion-panel-header class="light-green lighten-5">Базовые атрибуты</v-expansion-panel-header>
                <v-expansion-panel-content>

                 <div class="d-flex mt-4 align-center">
                   <h4 style="font-weight: 400;" class=" mr-4">Дата изготовления:</h4>
                   <date-picker
                       v-model="item['manufacture_date']" valueType="format"></date-picker>
                 </div>

                  <v-row>
                    <v-col cols="6" class="mb-0 pb-0">
                      <!-- Length -->
                      <v-text-field
                          class="mt-2"
                          label="Длина" dense
                          v-model.number="item['length']"
                      ></v-text-field>
                      <!-- Height -->
                      <v-text-field
                          class="mt-2"
                          label="Высота" dense
                          v-model.number="item['height']"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="6">
                      <!-- Width -->
                      <v-text-field
                          class="mt-2"
                          label="Ширина" dense
                          v-model.number="item['width']"
                      ></v-text-field>
                    </v-col>
                  </v-row>

                  <v-row class="mt-0 pt-0">
                    <v-col cols="6">
                      <!-- Weight -->
                      <v-text-field
                          class="mt-0"
                          label="Вес, гр." dense
                          v-model.number="item['weight']"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="6">
                      <!-- Volume -->
                      <v-text-field
                          class="mt-0"
                          label="Объем, мл." dense
                          v-model.number="item['volume']"
                      ></v-text-field>
                    </v-col>
                  </v-row>


                </v-expansion-panel-content>
              </v-expansion-panel>
            </v-expansion-panels>

            <v-switch
                class="mb-1 mt-6 mr-6 mt-0"
                :color="'teal'"
                persistent-hint
                messages="2"
                @change="updateItemData(true)"
                v-model="item['consider_weight']"
            >
              <template v-slot:label>
                  <span class="input__label w" :class="item['consider_weight'] ? 'teal--text' : 'grey--text text--darken-2'">
                    Весовой товар
                  </span>
              </template>
              <template v-slot:message>
                <span class="grey--text text--darken-2">Учет веса при расчетах (доставки, склада)</span>
              </template>
            </v-switch>

            <v-divider class="mb-4 mt-5"></v-divider>
            <!--  Целое или дробное количество товара -->

            <v-switch
                class="mb-2 mt-1 mr-6 mt-0"
                :color="'teal'"
                persistent-hint
                messages="2"
                @change="updateItemData(true)"
                v-model="item['is_integer']"
            >
              <template v-slot:label>
                  <span class="input__label w" :class="item['is_integer'] ? 'teal--text' : 'grey--text text&--darken-2'">
                    {{ item['is_integer'] ? 'Учет как целого' : 'Дробный учет' }}
                  </span>
              </template>
              <template v-slot:message>
                <span class="grey--text text--darken-2">Целый или дробный учет товара</span>
              </template>
            </v-switch>

          </v-col>

        </v-row>

        <v-divider class="mb-6 mt-8"></v-divider>

        <!-- Модель складского учета -->
        <v-data-table
            v-if="item"
            :headers="sourceItemsHeaders"
            :items.sync="item.source_items"
            class="elevation-1 mt-2 mb-4" style="border-top: 1px solid #ececec;"
            item-key="id"
        >
          <template v-slot:top>
            <v-toolbar color="white" flat class="mb-2">

              <v-toolbar-title>
                <h2 style="font-size: 21px;line-height: 1.1;">Модель складского учета<br>
                  <small style="font-weight: 400;font-size: 14px;">учет 1-ой складской ед. товара (SKU)</small>
                </h2>
              </v-toolbar-title>

              <v-divider
                  class="mx-8 mt-2"
                  inset
                  vertical
              ></v-divider>

              <v-switch
                  class="mb-2 mt-1 mr-6 mt-0"
                  :color="'teal'"
                  persistent-hint
                  messages="2"
                  @change="updateIsKit"
                  v-model="item['is_kit']"
              >
                <template v-slot:label>
                  <span class="input__label w" :class="item['is_kit'] ? 'teal--text' : 'grey--text text--darken-2'">
                    {{ item['is_kit'] ? 'Сборный товар' : 'Обычный товар' }}
                  </span>
                </template>
                <template v-slot:message>
                  <span class="grey--text text--darken-2">{{ item['is_kit'] ? 'Состоит из других товаров' : 'Состоит из себя же: 1 = 1' }}</span>
                </template>
              </v-switch>

              <v-divider
                  class="mx-8 mt-2"
                  inset
                  vertical
              ></v-divider>

              <div class="d-flex flex-column">
                <div class="lime--text text--darken-4 body-1">Стоимость товаров: <show-price :price="defaultRetailSum" color="grey--text text--darken-2"></show-price></div>
                <div class="light-green--text text--darken-3 body-1">Стоимость товаров со скидками: <show-price :price="defaultRetailDiscountSum"
                                                                                                        color="grey--text text--darken-2"></show-price></div>
              </div>

              <v-spacer></v-spacer>

              <v-card class="d-flex align-center" flat>

                <v-btn @click="uploadItem(true)"
                       color="secondary"
                       dark outlined
                       elevation="1"
                       small
                       class="align-center"
                >
                  <v-icon
                      class="ml-1 mr-3"
                      x-small
                  >far fa-sync
                  </v-icon>
                  Обновить
                </v-btn>

                <v-divider
                    class="mx-2 mt-0"
                    inset
                    vertical
                ></v-divider>

                <v-btn
                    @click="dialogChoiceProductSources = true"
                    outlined
                    :disabled="!item['is_kit']"
                    color="cyan darken-4"
                    elevation="1" small>
                  <v-icon
                      class="ml-1 mr-3"
                      x-small
                  >far fa-plus
                  </v-icon>
                  Добавить товар
                </v-btn>

              </v-card>

            </v-toolbar>
          </template>

          <template v-slot:item.source.retail_price="{ item }">
            <show-price :price="item.source['retail_price']"></show-price>
          </template>

          <template v-slot:item.source.retail_discount="{ item }">
            <show-price :price="item.source['retail_discount']"></show-price>
          </template>

          <!-- Сколько ед. в одном товаре -->
          <template v-slot:item.quantity="props" >
            <v-edit-dialog
                :return-value.sync="props.item.quantity"
                @save="saveAmountUnits(props.item)"
                :disabled="!isKit"
            >

              <div class="text-end item-editable-popup" style="min-width: 80px;">
                {{ props.item.source.is_integer ? (props.item.quantity ? props.item.quantity.toFixed(0) : '0') : (props.item.quantity ? props.item.quantity.toFixed(2) : '0') }}
              </div>

              <template v-slot:input >
                <v-currency-field
                    v-if="isKit"
                    :value-as-integer="props.item.source.is_integer"
                    :decimal-length="props.item.source.is_integer ? 0 : 2"
                    label="ед. в одном товаре"
                    :disabled="!isKit"
                    v-model.number="props.item.quantity">
                </v-currency-field>
              </template>

            </v-edit-dialog>
          </template>

          <template v-slot:item.source.measurement_unit.name="{ item }">
            <template v-if="item.source.measurement_unit_id">
              {{ item.source.measurement_unit.name }} ({{ item.source.measurement_unit.short_name }})
            </template>
          </template>

          <template v-slot:item.enable_viewing="{ item }">
            <show-item-yes-no class="item-editable-popup" :show="item['enable_viewing']" :bright="false" @click.native="changeEnableViewing(item)"
                              :disabled="!isKit"></show-item-yes-no>
          </template>

          <template v-slot:item.actions="{ item }">
            <show-action :route="{name:'product.edit', params: {public_id:item.source['public_id']}}" @delete="removeProductSource(item)"
                         :disabledRemove="!isKit"></show-action>
          </template>

        </v-data-table>

        <!-- Складской учет -->
        <transition name="fade">
          <v-data-table
            v-if="item && !item.is_kit"
            :headers="warehouseItemHeaders"
            :items.sync="item.warehouse_items"
            class="elevation-1 mt-2 mb-4" style="border-top: 1px solid #ececec;"
            item-key="id" sort-by="id"
        >
          <template v-slot:top>

            <v-toolbar color="white" flat>

              <v-toolbar-title><h2 style="font-size: 21px;">Складской учет</h2></v-toolbar-title>

              <v-spacer></v-spacer>

              <v-card class="d-flex align-center" flat>

                <!-- Обновление и сохранение серии -->
                <v-btn @click="uploadItem(true)" color="secondary" elevation="1" style="font-weight: 400;color:#101928;"
                       small outlined text>
                  <v-icon class="ml-1 mr-3" x-small>far fa-sync</v-icon>
                  Обновить
                </v-btn>
                <v-divider
                    class="mx-4 mt-2"
                    inset
                    vertical
                ></v-divider>
                <v-btn @click="dialogChoiceWarehouse = true"
                       color="secondary"
                       small outlined text
                       elevation="1" style="font-weight: 400;color:#101928;"
                       class="align-center"
                >
                  <v-icon class="ml-1 mr-3 green--text" x-small>far fa-plus</v-icon>
                  Добавить на склад
                </v-btn>

              </v-card>

            </v-toolbar>

          </template>

          <!-- Сколько ед. в одном товаре -->
          <template v-slot:item.stock="props" >
            <v-edit-dialog
                :return-value.sync="props.item.stock"
                @save="updateWarehouseItem(props.item)"
                :disabled="!isKit"
            >
              <div class="text-end item-editable-popup" style="min-width: 80px;">
                {{ isInteger ? (props.item.stock ? props.item.stock.toFixed(0) : '0') : (props.item.stock ? props.item.stock.toFixed(2) : '0') }}
              </div>

              <template v-slot:input >
                <v-currency-field
                    v-if="!isKit"
                    :value-as-integer="isInteger"
                    :decimal-length="isInteger ? 0 : 2"
                    label="ед. на складе"
                    v-model.number="props.item.stock">
                </v-currency-field>
              </template>

            </v-edit-dialog>
          </template>

          <template v-slot:item.reservation="props" >
            <v-edit-dialog
                :return-value.sync="props.item.reservation"
                @save="updateWarehouseItem(props.item)"
                :disabled="!isKit"
            >
              <div class="text-end item-editable-popup" style="min-width: 80px;">
                {{ isInteger ? (props.item.reservation ? props.item.reservation.toFixed(0) : '0') : (props.item.reservation ? props.item.reservation.toFixed(2) : '0') }}
              </div>

              <template v-slot:input >
                <v-currency-field
                    v-if="!isKit"
                    :value-as-integer="isInteger"
                    :decimal-length="isInteger ? 0 : 2"
                    label="ед. в резерве"
                    v-model.number="props.item.reservation">
                </v-currency-field>
              </template>

            </v-edit-dialog>
          </template>

          <!--<template v-slot:item.reservation="{ item }">
            <span style="color: #794f09;">{{item.reservation}}</span>
          </template>-->

          <template v-slot:item.all_in_stock="{ item }">
            <span style="color: #415663;">{{item.stock+item.reservation}}</span>
          </template>

          <template v-slot:item.actions="{ item }">
            <show-action :route="{name:'product.edit', params: {public_id:item.product['public_id']}}" @delete="removeProduct(item)" removeIcon></show-action>
          </template>

        </v-data-table>
        </transition>

        <!-- Список карточек товара-->
        <v-data-table
              v-if="item"
              :headers="[
              { text: 'Id', align: 'start', value: 'public_id'},
              { text: 'Наименование', value: 'label'},
              { text: 'Доп. наименование', value: 'second_label'},
              { text: 'Теги', value: 'product_tags', sortable: false},
              { text: 'Path',  value: 'path'},
              { text: 'Route name',  value: 'route_name'},
              { text: 'Товаров', align: 'end', value: '_product_count'},
              { text: 'Статус розн.', align: 'end', value: 'enable_retail_sale'},
              { text: 'Статус опт.',  align: 'end',value: 'enable_wholesale_sale'},
              { text: 'Действие', align: 'center', value: 'actions', sortable: false},
            ]"
              :items.sync="item.product_cards"
              class="elevation-1 mt-2 mb-4" style="border-top: 1px solid #ececec;"
              item-key="id"
          >
          <template v-slot:top>
            <v-toolbar color="white" flat>
              <v-toolbar-title><h2 style="font-size: 21px;">Список карточек товаров</h2></v-toolbar-title>
              <v-spacer></v-spacer>
              <v-card class="d-flex align-center" flat>

                <v-btn @click="uploadItem(true)"
                       color="secondary"
                       dark outlined
                       elevation="1"
                       small
                       class="align-center"
                >
                  <v-icon
                      class="ml-1 mr-3"
                      x-small
                  >far fa-sync</v-icon>
                  Обновить
                </v-btn>

                <v-divider
                    class="mx-2 mt-0"
                    inset
                    vertical
                ></v-divider>

                <v-btn
                    @click="dialogChoiceProductCard = true"
                    outlined
                    elevation="1"
                    small
                >
                  <v-icon
                      class="cyan--text ml-1 mr-3"
                      x-small
                  >far fa-plus</v-icon>
                  Добавить карточку
                </v-btn>

              </v-card>

            </v-toolbar>
          </template>

          <template v-slot:item.enable_retail_sale="{ item }">
            <show-item-yes-no :show="item['enable_retail_sale']"></show-item-yes-no>
          </template>

          <template v-slot:item.enable_wholesale_sale="{ item }">
            <show-item-yes-no :show="item['enable_wholesale_sale']"></show-item-yes-no>
          </template>

          <template v-slot:item.product_tags="{ item }">
            <show-tags :tags="item.product_tags"></show-tags>
          </template>

          <template v-slot:item.product_count="{ item }">
            <span v-text="item['products'] && item['products'].length > 0 ? item['products'].length : '0'" :class="{ 'red--text': !(item['products'] && item['products'].length > 0) }"></span>
          </template>

          <template v-slot:item.actions="{ item }">
            <show-action :route="{name:'product-cards.edit', params: {public_id:item['public_id']}}" @delete="removeProductFromProductCard(item)" removeIcon></show-action>
          </template>

        </v-data-table>

        <!--  Складской учет . Изображения / Атрибуты / Описание  -->
        <v-expansion-panels class="px-0 mx-0 mt-3 text--white" multiple hover tile>
          <v-expansion-panel style="background:#f1f2e7;color:white">
            <v-expansion-panel-header class="body-1 grey--text text--darken-4">Изображения товара ({{ item.images.length }} шт.)</v-expansion-panel-header>
            <v-expansion-panel-content>
              <div class="d-flex justify-end">
                <div class="d-flex flex-column" style="min-width: 300px;">
                  <v-file-input
                      @change.native="uploadImage"
                      dense
                      style="max-width: 300px" class="mb-2"
                      label="Добавить файл"
                      show-size
                      v-model="uploadedFile" counter></v-file-input>
                </div>
              </div>
              <draggable
                  :list="item.images"
                  @end="savePriority"
              >
                <transition-group name="flip-list">
                  <v-card v-for="image in item.images" :key="image.id" class="mt-1 pl-1">
                    <v-row>
                      <v-col cols="3">
                        <a :href="image['_url']" target="_blank" style="width: 100px;">
                          <v-img :src="image['_url']" aspect-ratio="1" max-width="100px"></v-img>
                        </a>
                      </v-col>
                      <v-col cols="9">
                        <v-text-field v-model="image.alt" placeholder="укажите alt" dense></v-text-field>
                        <section class="d-flex justify-space-between">
                          <span class="body-2 blue-grey--text">{{ formatBytes(image['size']) }}</span>
                          <span class="body-2 grey--text">priority: {{ image['priority'] }}</span>
                          <div>
                            <v-icon @click="updateImage(image)" small class="blue--text text--lighten-2 mr-4">fas fa-save</v-icon>
                            <v-icon @click="deleteImage(image)" small class="red--text text--lighten-2 mr-2">fas fa-trash</v-icon>
                          </div>
                        </section>
                      </v-col>
                    </v-row>
                  </v-card>
                </transition-group>
              </draggable>
            </v-expansion-panel-content>
          </v-expansion-panel>
          <v-expansion-panel style="background:#daece8;color:white">
            <v-expansion-panel-header class="body-1 grey--text text--darken-4">Атрибуты товара</v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-card flat class="pa-2 ma-0">
                <v-row>
                  <v-col cols="12">
                    <v-spacer></v-spacer>
                    <v-card class="d-flex justify-end mb-3" flat style="background-color: transparent;">
                      <v-btn @click="newProductAttribute"
                             color="secondary" dark small outlined>
                        <v-icon class="ml-1 mr-3" x-small>far fa-plus</v-icon>
                        Добавить свойство
                      </v-btn>
                      <v-divider class="mx-2 mt-0" inset vertical></v-divider>
                      <v-btn @click="updateItemData()"
                             color="teal darken-3" dark small outlined>
                        <v-icon class="ml-1 mr-3" x-small>far fa-save</v-icon>
                        Сохранить
                      </v-btn>
                    </v-card>
                  </v-col>
                  <v-col v-for="(_,key, index) in attributes" :key="index" cols="6">
                    <v-row style="border: 1px dotted grey;" class="pa-1" dense>
                      <v-col cols="6" class="pr-1">
                        <v-text-field :value="key" label="Key" @change="updateKeyNameAttribute(item, key, $event, true)"></v-text-field>
                      </v-col>
                      <v-col cols="6" class="pl-1">
                        <v-text-field v-model="item.attributes[key]" label="value" append-icon="far fa-trash" @click:append="deleteAttribute(item, key)"></v-text-field>
                      </v-col>
                    </v-row>

                  </v-col>
                </v-row>
              </v-card>
            </v-expansion-panel-content>
          </v-expansion-panel>
          <v-expansion-panel style="background:#eceada;color:white">
            <v-expansion-panel-header class="body-1 grey--text text--darken-4">Описание товара</v-expansion-panel-header>
            <v-expansion-panel-content>
              <div class="d-flex justify-end mb-2">
                <v-btn small class="mb-1" @click="item.description='<p>Описание карточки товара..</p>'">Set HTML tpl</v-btn>
              </div>
              <section class="pb-6">
                <prism-editor v-model="item.description" class="myPrismEditor" language="html"></prism-editor>
              </section>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>

        <v-divider class="mb-6 mt-8"></v-divider>

      </v-form>
    </v-card>

    <choice-manufacturer :open.sync="dialogChoiceManufacturer" @choice="choiceManufacturer" @close="dialogChoiceManufacturer = false"></choice-manufacturer>
    <choice-categories v-if="item" :open="dialogChoiceCategories" @choice="choiceCategories" @close="dialogChoiceCategories = false"
                       :catIds.sync="item.product_categories"></choice-categories>
    <choice-product-tags v-if="item" :open="dialogChoiceTags" @choice="choiceTags" @close="dialogChoiceTags = false"
                         :catIds.sync="item.product_tags"></choice-product-tags>
    <choice-product-type v-if="item" :open="dialogChoiceType" @choice="choiceType" @close="dialogChoiceType = false"
                         :choiceId.sync="item.product_type_id"></choice-product-type>
    <choice-product :open.sync="dialogChoiceProductSources" @choice="choiceProductSource" @close="dialogChoiceProductSources = false"
                    :isKit="item && item.is_kit ? false : null"></choice-product>
    <choice-product-card :open.sync="dialogChoiceProductCard" @choice="appendProductToProductCard" @close="dialogChoiceProductCard = false"></choice-product-card>
    <choice-warehouse :open.sync="dialogChoiceWarehouse" @choice="appendProductToWarehouse" @close="dialogChoiceWarehouse = false" :wIdx.sync="warehouseIdx"></choice-warehouse>

  </v-container>
</template>

<script>
import draggable from 'vuedraggable'
import _ from 'lodash'
const  preloadsProduct = "Images,ProductCards,ProductType,ProductCategories,Manufacturer,MeasurementUnit,VatCode,ProductTags,SourceItems.Source.MeasurementUnit," +
    "WarehouseItems,WarehouseItems.Warehouse"
export default {
  data: () => ({
    formValid: true, // form
    item: null,
    uploadedFile: null, // загружаемый файл изображения
    tab: null,
    enabledEdit: true,

    measurementUnits: [],
    vatCodes: [],
    paymentSubjects: [],

    loading: false, // абстрактный статус загрузки..
    dialogChoiceManufacturer: false,
    dialogChoiceCategories: false, // dialog template
    dialogChoiceType: false, // dialog template
    dialogChoiceTags: false, // dialog template
    dialogChoiceProductSources: false,
    dialogChoiceProductCard: false,
    dialogChoiceWarehouse: false,

    errors: [],
    sourceItemsHeaders: [
      {text: 'Id', align: 'start', value: 'source.public_id'},
      {text: 'Наименование', value: 'source.label'},
      {text: 'Артикул', value: 'source.article'},
      {text: 'Модель', value: 'source.model'},
      {text: 'Розничная цена', value: 'source.retail_price', align: 'end'},
      {text: 'Скидка', value: 'source.retail_discount', align: 'end'},
      {text: 'Доступен просмотр', value: 'enable_viewing', align: 'end'},
      {text: 'Ед. в одном товаре', value: 'quantity', align: 'end'},
      {text: 'Ед. измерения', value: 'source.measurement_unit.name', align: 'end'},
      {text: 'Действия', value: 'actions', sortable: false, align: 'center'},
    ],
    warehouseItemHeaders: [
      { text: 'Id', align: 'start', value: 'id'},
      { text: 'Склад', value: 'warehouse.label', sortable: false},
      { text: 'Краткое', value: 'warehouse.short_label', sortable: false},
      { text: 'Код', value: 'warehouse.code', sortable: false},
      // { text: 'Время хранения', align: 'center', value: 'expired_at'},
      { text: 'Доступно к заказу', align: 'end', value: 'stock'},
      { text: 'В резерве', align: 'end', value: 'reservation'},
      { text: 'Всего', align: 'end', value: 'all_in_stock', sortable: false},

      { text: 'Действие', align: 'center', value: 'actions', sortable: false},
    ],
    warehouseIdx: [], // Список складов
  }),
  async mounted() {
    await Promise.all([
      this.getMeasurementUnitsFromApi(),
      this.getVatCodesFromApi(),
      this.getPaymentSubjectsFromApi(),
    ]).then(() => this.uploadItem())
  },
  methods: {
    uploadItem: async function (showNotification = false) {
      this.loading = true
      await this.$api.product.getByPublicId({
        accountHashId: this.$store.state.account.hash_id,
        public_id: this.$route.params.public_id,
        preloads: preloadsProduct
      })
          .then(resp => {
            if (resp['product'] !== undefined) {
              this.item = resp['product']

              if (resp['product']['images']) this.item.images = this.sortByKey(resp['product']['images'], 'priority')
              if (!resp['product']['images']) this.item['images'] = []

              if (!resp['product']['attributes']) this.item['attributes'] = {}

              // Список складов..
              this.warehouseIdx = _.map(this.item.warehouse_items,'warehouse_id')

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
              title: 'Кампания не найдена',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)

    },
    async uploadImage() {
      await this.uploadFile(this.item.id, "products")
          .then((resp) => {
            if (resp['file'] !== undefined && resp['file'] !== null) {
              this.uploadedFile = null
              if (!this.item.images) this.item["images"] = []
              this.item.images.unshift(resp['file'])
              this.$notify({
                group: 'main',
                title: 'Файл создан',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Не полные данные от сервера',
                type: 'warring',
              });
            }
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка при создании файла',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
            this.uploadedFile = null
          });
    },
    async updateImage(image) {
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: image.id,
      }
      Object.assign(payload, image)
      this.$api.storage.update(payload)
          .then((resp) => {

            if (resp['file']) {
              let index = this.item.images.findIndex(el => el.id === image.id)
              if (index !== -1) this.item.images[index] = resp['file']
            }

            this.$notify({
              group: 'main',
              title: 'Изображение сохранено',
              type: 'success',
            });
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка сохранения',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => {
          })
    },
    async deleteImage(image) {
      const payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: image.id,
      }
      this.$api.storage.delete(payload)
          .then(() => {
            let index = this.item.images.findIndex(el => el.id === image.id)
            if (index !== -1) {
              this.item.images.splice(index, 1);
            }
            this.item.image = null
            this.$notify({
              group: 'main',
              title: 'Изображение удалено',
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
          .finally(() => {
          })
    },
    updateItemData: async function (showNotification = false) {
      let payload = this.item

      // if (payload['volume'].length < 1) payload['volume'] = null

      Object.assign(payload, {
        accountHashId: this.$store.state.account.hash_id,
        preloads: preloadsProduct
      })

      await this.$api.product.update(payload)
          .then((resp) => {
            if (resp['product'] !== undefined) {
              this.item = resp['product']
              if (resp['product']['images']) this.item.images = this.sortByKey(resp['product']['images'], 'priority')
              if (showNotification)
                this.$notify({
                  group: 'main',
                  title: 'Данные обновлены',
                  type: 'success',
                });
            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [product]',
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
    choiceManufacturer: async function (manufacturer) {

      if (!manufacturer || !manufacturer.id) return

      let payload = {
        id: this.item.id,
        public_id: this.item.public_id,
        manufacturer_id: manufacturer.id,
        accountHashId: this.$store.state.account.hash_id,
        preloads: "Manufacturer",
      }

      await this.$api.product.update(payload)
          .then((resp) => {
            if (resp['product'] !== undefined) {
              this.item.manufacturer = resp['product']['manufacturer']
              this.item.manufacturer_id = resp['product']['manufacturer']['id']

              this.$notify({
                group: 'main',
                title: 'Производитель обновлен',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [product]',
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
          .finally(() => this.dialogChoiceManufacturer = false)

    },
    deleteManufacturer: async function () {
      this.item.manufacturer_id = null
      await this.updateItemData()
    },
    deleteType: async function () {
      this.item.product_type_id = null
      await this.updateItemData()
    },
    openDialogChoiceManufacturer() {
      this.dialogChoiceProduct = true
    },

    savePriority: async function (showNotification = false, force = false) {

      let arr = []

      this.item.images.forEach((item, i) => {
        if (force) {
          let v = {
            id: item.id,
            priority: i + 1,
            owner_id: item.id,
            owner_type: item.owner_type,
          }
          let index = arr.findIndex(el => el.id === item.id)
          if (index !== -1) {
            arr[index].priority = (i + 1)
          } else {
            arr.push(v)
          }
        } else if (item.priority !== (i + 1)) {
          let v = {
            id: item.id,
            priority: i + 1,
            owner_id: this.item.id,
            owner_type: item.owner_type,
          }
          let index = arr.findIndex(el => el.id === item.id)
          if (index !== -1) {
            arr[index].priority = (i + 1)
          } else {
            arr.push(v)
          }
        }
      })

      if (arr.length < 1) {return}

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        files: arr,
      }

      this.$api.storage.massUpdates(payload)
          .then(resp => {
            if (resp['files'] !== undefined) {
              this.item.images = this.sortByKey(resp['files'], 'priority')
              this.carried = false
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
              title: 'Ошибка обновления порядка',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.loading = false)
    },

    appendProductToProductCard: async function (productCard) {

      if (!productCard || !productCard.id) return

      let payload = {
        id: this.item.id,
        accountHashId: this.$store.state.account.hash_id,
        product_card_id: productCard.id,
        preloads: preloadsProduct,
      }

      this.$api.product.appendToProductCard(payload)
          .then(async (resp) => {
            if (resp['product']) {

              if (!resp['product']['product_cards']) this.item['product_cards'] = []
              else this.item['product_cards'] = resp['product']['product_cards']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе product',
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
          .finally(() => {
            this.dialogChoiceProductCard = false
          })

    },

    removeProductFromProductCard: async function (productCard) {

      let payload = {
        id: this.item.id,
        accountHashId: this.$store.state.account.hash_id,
        product_card_id: productCard.id,
        preloads: preloadsProduct,
      }

      this.$api.product.removeFromProductCard(payload)
          .then(async (resp) => {
            if (resp['product']) {

              if (!resp['product']['product_cards']) this.item['product_cards'] = []
              else this.item['product_cards'] = resp['product']['product_cards']

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе product',
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

    },

    getMeasurementUnitsFromApi: async function (showNotification = false) {

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
      }

      await this.$api.measurementUnit.getList(payload)
          .then(resp => {
            if (resp['measurement_units'] !== undefined) {
              this.measurementUnits = resp['measurement_units']
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка получения списка measurement units',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка получения списка Payment Subjects',
              text: err['message'],
              type: 'error',
            });
          })
    },
    getVatCodesFromApi: async function (showNotification = false) {

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
      }

      await this.$api.vatCode.getList(payload)
          .then(resp => {
            if (resp['vat_codes'] !== undefined) {
              this.vatCodes = resp['vat_codes']
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка получения списка vat codes',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка получения списка Vat Codes',
              text: err['message'],
              type: 'error',
            });
          })
    },
    getPaymentSubjectsFromApi: async function (showNotification = false) {

      const payload = {
        accountHashId: this.$store.state.account.hash_id,
      }

      await this.$api.paymentSubject.getList(payload)
          .then(resp => {
            if (resp['payment_subjects'] !== undefined) {
              this.paymentSubjects = resp['payment_subjects']
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка получения списка payment Subjects',
                type: 'warring',
              });
            }

          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка получения списка Payment Subjects',
              text: err['message'],
              type: 'error',
            });
          })
    },

    choiceCategories: async function (data) {

      let categories = []
      let ids = []

      if (data.append) {
        for (let _index in this.item['product_categories']) {
          if (this.item['product_categories'].hasOwnProperty(_index) && typeof this.item['product_categories'][_index] === 'object') {
            ids.push(this.item['product_categories'][_index].id)
          } else {
            ids.push(this.item['product_categories'][_index])
          }
        }

        categories = _.union(ids, data.categories);

      } else {
        // надо подготовить список idx
        for (let _index in this.item['product_categories']) {
          if (this.item['product_categories'].hasOwnProperty(_index) && typeof this.item['product_categories'][_index] === 'object') {
            ids.push(this.item['product_categories'][_index].id)
          } else {
            ids.push(this.item['product_categories'][_index])
          }
        }

        categories = _.difference(ids, data.categories);

      }


      let arr = []
      for (let index in categories) {
        arr.push({id: categories[index]})
      }

      // Обновляем данные
      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_categories: arr,
        preloads: "ProductCategories",
      }

      await this.$api.product.syncProductCategories(payload)
          .then((resp) => {
            if (resp['product'] !== undefined) {
              this.item['product_categories'] = resp['product']['product_categories']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [email_notification]',
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
          .finally(() => this.dialogChoiceCategories = false)
    },
    choiceTags: async function (data) {

      let tags = []
      let ids = []

      if (data.append) {
        for (let _index in this.item['product_tags']) {
          if (this.item['product_tags'].hasOwnProperty(_index) && typeof this.item['product_tags'][_index] === 'object') {
            ids.push(this.item['product_tags'][_index].id)
          } else {
            ids.push(this.item['product_tags'][_index])
          }
        }

        tags = _.union(ids, data.categories);

      } else {
        // надо подготовить список idx
        for (let _index in this.item['product_tags']) {
          if (this.item['product_tags'].hasOwnProperty(_index) && typeof this.item['product_tags'][_index] === 'object') {
            ids.push(this.item['product_tags'][_index].id)
          } else {
            ids.push(this.item['product_tags'][_index])
          }
        }

        tags = _.difference(ids, data.categories);

      }


      let arr = []
      for (let index in tags) {
        arr.push({id: tags[index]})
      }

      // Обновляем данные
      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_tags: arr,
        preloads: "ProductTags",
      }

      await this.$api.product.syncProductTags(payload)
          .then((resp) => {
            if (resp['product'] !== undefined) {
              this.item['product_tags'] = resp['product']['product_tags']
              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Данные от сервера не полные',
                text: 'не хватает объекта [product_tags]',
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
          .finally(() => this.dialogChoiceTags = false)
    },
    choiceType: async function (productType) {

      this.item.product_type_id = parseInt(productType.id)

      await this.updateItemData(false)
          .finally(() => this.dialogChoiceType = false)
    },

    deleteAttribute: async function (item, keyName) {
      delete item['attributes'][keyName]

      await this.updateItemData()
    },
    insertProductAttribute: async function (keyName, value) {
      if (!this.item['attributes']) this.item['attributes'] = {}
      this.$set(this.item['attributes'], keyName, value)

      await this.updateItemData()

    },
    updateKeyNameAttribute: _.throttle(async function (item, oldName, newName, showNotification) {
      item['attributes'][newName] = item['attributes'][oldName]
      delete item['attributes'][oldName]

      await this.updateItemData()
    }, 1400),
    newProductAttribute: async function () {
      let attrName = "new_key_"
      for (let index = 0; index < 100; index++) {
        if (!this.attributes.hasOwnProperty(attrName + index)) {
          attrName += index
          break
        }
      }
      // this.insertProductAttribute(item, 'new_key_' + Math.round(Math.random()) + '_' + Math.round(Math.random()), '-')
      await this.insertProductAttribute(attrName, '')

    },
    saveAmountUnits: async function(item){

      if (!this.item['is_kit']) return

      await this.updateSourceItem(item.source_id)
    },
    choiceProductSource: async function (sourceProduct) {

      if (!sourceProduct || !sourceProduct.id) return


      if (this.isKit && sourceProduct.id === this.item.id) {
        this.$notify({
          group: 'main',
          title: 'Ошибка добавление товара',
          text: 'Сборный товар не может состоять из себя же',
          type: 'warring',
        });
        return
      }

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        source_id: sourceProduct.id,
        quantity: 0, // ед. в одном товаре
        enable_viewing: true, // отображать в списке
        preloads: preloadsProduct
      }

      this.$api.product.appendProductSource(payload)
          .then(async (resp) => {
            if (resp['product']) {
              if (resp['product']['source_items'])
                this.item.source_items = resp['product']['source_items']
              else
                this.item.source_items = []

              this.$notify({
                group: 'main',
                title: 'Данные обновлены',
                type: 'success',
              });
            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе source_items',
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
          .finally(async () => {
            // await this.getTagDataFromApi()
            this.dialogChoiceProductSources = false
          })

    },
    removeProductSource: async function (sourceItem) {

      if (!sourceItem || !sourceItem.source_id) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        source_id: sourceItem.source_id,
        preloads: preloadsProduct
      }

      this.$api.product.removeProductSource(payload)
        .then(async (resp) => {
          if (resp['product']) {
            if (resp['product']['source_items'])
              this.item.source_items = resp['product']['source_items']
            else
              this.item.source_items = []

            this.$notify({
              group: 'main',
              title: 'Данные обновлены',
              type: 'success',
            });
          } else {
            this.$notify({
              group: 'main',
              title: 'Ошибка удаления',
              text: 'Отсутствует переменная в ответе source_items',
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
          .finally(async () => {
            // await this.getTagDataFromApi()
            this.dialogChoiceProductSources = false
          })

    },

    // Обновляет по sourceId
    updateSourceItem(sourceId) {

      let _item = this.item.source_items.find(el => el.source_id === sourceId)
      if (!_item) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          text: 'Объект source item не найден',
          type: 'error',
        });
        return
      }

      // Проверяем чтобы Кит не состоял из себя же
      if (this.isKit && this.item.id === sourceId) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          text: 'Сборный товар не может состоять из себя же',
          type: 'error',
        });
        return
      }

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        source_id: sourceId,
        quantity: _item.quantity, // ед. в одном товаре
        enable_viewing: _item.enable_viewing,
        preloads: preloadsProduct
      }

      this.$api.product.updateProductSource(payload)
          .then(async (resp) => {
            if (resp['product']) {

              if (resp['product']['source_items'])
                this.item.source_items = resp['product']['source_items']
              else
                this.item.source_items = []

              this.$notify({
                group: 'main',
                title: 'Данные модели товара обновлены',
                type: 'success',
              });

            } else {
              this.$notify({
                group: 'main',
                title: 'Ошибка обновления',
                text: 'Отсутствует переменная в ответе source_items',
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
    changeEnableViewing: async function (item) {

      if (!this.item['is_kit']) return

      item.enable_viewing = !item.enable_viewing

      await this.updateSourceItem(item.source_id)

    },
    updateIsKit: async function() {
      // Проверки
      if (!this.item || !this.item.source_items) {
        this.$notify({
          group: 'main',
          title: 'Ошибка обновления',
          text: 'Объект source item не найден',
          type: 'error',
        });
        this.item.is_kit = !this.item.is_kit
        return
      }

      // Собираем массив замены
      let productSources = []
      if (!this.item.is_kit) {

        productSources = [{
          product_id: this.item.id,
          source_id: this.item.id,
          quantity: 1.0,
          enable_viewing: true
        }]
      }

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: this.item.id,
        product_sources: productSources,
        preloads: preloadsProduct
      }

      this.updateItemData().then(()=>{
        this.$api.product.syncProductSource(payload)
            .then(async (resp) => {
              if (resp['product']) {

                if (resp['product']['source_items'])
                  this.item.source_items = resp['product']['source_items']
                else
                  this.item.source_items = []

              } else {
                this.$notify({
                  group: 'main',
                  title: 'Ошибка обновления',
                  text: 'Отсутствует переменная в ответе source_items',
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
      })

    },

    removeProduct: async function (warehouseItem) {
      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: warehouseItem.id,
      }

      this.$api.warehouseItem.delete(payload)
          .then(async () => this.uploadItem())
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка создания',
              text: err['message'],
              type: 'error',
            });
          })

    },
    appendProductToWarehouse: async function (warehouse) {

      if (!warehouse || !warehouse.id) return

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: warehouse.id,
        product_id: this.item.id,
      }

      this.$api.warehouse.appendProduct(payload)
          .then(async () => this.uploadItem())
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка создания',
              text: err['message'],
              type: 'error',
            });
          })
          .finally(() => this.dialogChoiceWarehouse = false)

    },
    updateWarehouseItem: async function(warehouseItem) {
      console.log(warehouseItem)

      let payload = {
        accountHashId: this.$store.state.account.hash_id,
        id: warehouseItem.id,
        product_id: warehouseItem.product_id,
        stock: warehouseItem.stock,
        reservation: warehouseItem.reservation,

      }
      this.$api.warehouseItem.update(payload)
          .then(async () => {
            this.$notify({
              group: 'main',
              title: 'Данные на складе обновлены',
              type: 'success',
            });
            await this.uploadItem()
          })
          .catch((err) => {
            this.$notify({
              group: 'main',
              title: 'Ошибка обновления данных склада',
              text: err['message'],
              type: 'error',
            });
          })
    },

  },
  computed: {
    attributes() {
      if (!this.item) return {}
      let arr = {}
      Object.keys(this.item['attributes']).sort().forEach((key, index) => {
        arr[key] = this.item['attributes'][key]
      })

      return arr
    },
    isKit() {
      return this.item.is_kit
    },
    isInteger() {
      return this.item.is_integer
    },
    defaultRetailSum(){
      return _.sumBy(this.item.source_items,(el)=>{return el.quantity * el.source.retail_price} )
    },
    defaultRetailDiscountSum(){
      return _.sumBy(this.item.source_items,(el)=>{return el.quantity * (el.source.retail_price-el.source.retail_discount)} )
    }
  },
  components: {
    draggable,
    DelayTime: () => import('@/components/views/tpl/email-marketing/queue/layouts/DelayTime'),
    ChoiceManufacturer: () => import('@/components/views/tpl/organizations/manufacturer/ChoiceManufacturer'),
    ChoiceCategories: () => import('@/components/views/tpl/shops/product-category/ChoiceCategories'),
    ChoiceProductTags: () => import('@/components/views/tpl/shops/product-tag/ChoiceTags'),
    ChoiceProductType: () => import('@/components/views/tpl/shops/product-type/ChoiceProductType'),
    ChoiceProductCard: () => import('@/components/views/tpl/shops/product-cards/ChoiceProductCard'),
    ChoiceWarehouse: () => import('@/components/views/tpl/warehouse/warehouse/ChoiceWarehouse'),
    ShowPrice: () => import('@/components/views/tpl/layouts/table/ShowPrice'),
    ShowAction: () => import('@/components/views/tpl/layouts/table/ShowAction'),
    ShowTags: () => import('@/components/views/tpl/layouts/table/ShowTags'),
    ChoiceProduct: () => import('@/components/views/tpl/shops/products/ChoiceProduct'),
    ShowItemYesNo: () => import('@/components/views/tpl/layouts/table/ShowItemYesNo'),
  },
}
</script>

