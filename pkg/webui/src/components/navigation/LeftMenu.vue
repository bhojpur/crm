<template>
  <v-navigation-drawer
      :expand-on-hover="false"
      :mini-variant="miniVariant"
      app
      class="navigationDrawerClass"
      clipped
      color="#162027"
      dark
      disable-resize-watcher
      fixed
      permanent
      style="border-right: none !important;padding-top: 15px;"

      width="240"
  >
    <v-list dense>
      <v-list-group
          :key="item.title"
          :to="{name: item.route_name}"
          active-class="nav-active-item"
          append-icon="far fa-angle-down"
          dark
          exact
          v-for="item in this.items"
          v-if="item['items']!== undefined && item['items'].length > 0"
          v-model="item.active"
      >
        <template v-slot:activator>
          <v-list-item-icon class="mr-4 justify-center">
            <v-icon small v-text="item.icon"/>
          </v-list-item-icon>
          <v-list-item-content>
            <v-list-item-title dark style="font-size: 14px;font-weight: normal;" v-text="item.title"/>
          </v-list-item-content>
        </template>

        <v-list-item
            :key="subItem.title"
            :to="{name: subItem.route_name}"
            dark
            exact-active-class="exact-nav-list-item--active"
            v-for="subItem in item.items"

        >
          <v-list-item-icon class="ml-4 mr-3 justify-center">
            <v-icon small v-text="subItem.icon"/>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title style="font-size: 13px;" v-text="subItem.title"></v-list-item-title>
          </v-list-item-content>

        </v-list-item>

      </v-list-group>

      <v-list-item
          :key="item.title"
          :to="{name: item.route_name}"
          active-class="nav-active-item"
          dark
          exact
          exact-active-class="exact-nav-list-item--active"
          v-else
      >
        <v-list-item-icon class="mr-4 justify-center">
          <v-icon small v-text="item.icon"/>
        </v-list-item-icon>

        <v-list-item-content>
          <v-list-item-title style="font-size: 14px;" v-text="item.title"></v-list-item-title>
        </v-list-item-content>

      </v-list-item>

    </v-list>
    <v-spacer></v-spacer>
<!--    <v-switch v-model="$vuetify.theme.dark" label="theme"></v-switch>-->

  </v-navigation-drawer>
</template>

<script>
export default {
  data: () => ({
    drawer: true,
    items: [
      {
        title: 'Home',
        icon: 'far fa-home',
        route_name: 'home',
        active: true,
      },

      {
        title: 'Заказы',
        icon: 'far fa-cart-plus',
        route_name: 'orders.list',
        active: false,
      },
      {
        title: 'Вопросы',
        icon: 'far fa-map-marker-question',
        route_name: 'questions.list',
        active: false,
      },
      {
        title: 'Доставка товаров',
        icon: 'fal fa-shipping-fast',
        route_name: 'delivery.orders',
        active: false,
      },
      {
        title: 'Номенклатура',
        icon: 'fal fa-store-alt',
        route_name: 'shops',
        active: false,
        items: [
          {
            title: 'Товары',
            // icon: 'fal fa-gifts',
            icon: 'far fa-cubes',
            route_name: 'products.list',
            active: false,
          },
          {
            title: 'Карточки',
            // icon: 'fal fa-layer-group',
            icon: 'fal fa-th-list',
            route_name: 'product-cards.list',
            active: false,
          },
          {
            title: 'Категории',
            // icon: 'fal fa-layer-group',
            icon: 'fal fa-burger-soda',
            route_name: 'product-categories.list',
            active: false,
          },
          {
            title: 'Типы товаров',
            // icon: 'fal fa-layer-group',
            icon: 'fal fa-layer-group',
            route_name: 'product-types.list',
            active: false,
          },
          {
            title: 'Теги',
            // icon: 'fal fa-layer-group',
            icon: 'fal fa-tag',
            route_name: 'product-tags.list',
            active: false,
          },
          {
            title: 'Группы тегов',
            // icon: 'fal fa-layer-group',
            icon: 'fal fa-tags',
            route_name: 'product-tag-groups.list',
            active: false,
          },
        ]
      },
      {
        title: 'Складской учет',
        icon: 'fal fa-boxes-alt',
        route_name: 'inventory-control',
        active: false,
        items: [
          {
            title: 'Поставки',
            // icon: 'fal fa-layer-group',
            icon: 'fal fa-th-list',
            route_name: 'shipments.list',
            active: false,
          },
          {
            title: 'Инвентаризации',
            // icon: 'fal fa-layer-group',
            icon: 'fal fa-th-list',
            route_name: 'inventories.list',
            active: false,
          },
          {
            title: 'Склады',
            // icon: 'fal fa-gifts',
            icon: 'far fa-cubes',
            route_name: 'warehouses.list',
            active: false,
          },

        ]
      },
      {
        title: 'Контакты',
        icon: 'fal fa-book-user',
        route_name: 'organizations.*',
        active: false,
        items: [
          {
            title: 'Персоналии',
            icon: 'fal fa-user-friends',
            route_name: 'users.list',
            active: false,
          },
          {
            title: 'Производители',
            icon: 'fal fa-digging',
            route_name: 'manufacturers.list',
            active: false,
          },
          {
            title: 'Контрагенты',
            icon: 'fal fa-industry-alt',
            route_name: 'companies.list',
            active: false,
          },
        ]
      },
      {
        title: 'Web-сайты',
        icon: 'fal fa-space-station-moon-alt',
        route_name: 'sites',
        active: false,
        items: [
          {
            title: 'Страницы',
            icon: 'fal fa-folder-tree',
            route_name: 'web-pages.list',
            active: false,
          },
          {
            title: 'Сайты',
            icon: 'fal fa-store-alt',
            route_name: 'web-sites.list',
            active: false,
          },

        ]
      },
      {
        title: 'Контент',
        icon: 'fal fa-server',
        route_name: 'content.storage',
        active: false,
        items: [
          {
            title: 'Файлы',
            icon: 'fal fa-folder-tree',
            route_name: 'content.storage',
            active: false,
          },
          {
            title: 'Статьи',
            icon: 'fal fa-newspaper',
            route_name: 'content.articles',
            active: false,
          },
        ]
      },
      {
        title: 'Email-маркетинг',
        icon: 'fal fa-mail-bulk',
        route_name: 'marketing.email',
        active: false,
        items: [
          {
            title: 'Уведомления',
            icon: 'fal fa-envelope-open-text',
            route_name: 'email-marketing.notifications',
            active: false,
          },
          {
            title: 'Массовые рассылки',
            icon: 'fal fa-mail-bulk',
            route_name: 'email-marketing.campaigns',
            active: false,
          },
          {
            title: 'Автоматические серии',
            icon: 'fal fa-send-back',
            route_name: 'email-marketing.queues',
            active: false,
          },

          {
            title: 'Шаблоны',
            icon: 'fal fa-line-columns',
            route_name: 'email-marketing.templates',
            active: false,
          },
        ]
      },
      {
        title: 'Финансы',
        icon: 'fal fa-hand-holding-usd',
        route_name: 'finance.*',
        active: false,
        items: [
          {
            title: 'Платежи',
            icon: 'fal fa-file-invoice-dollar',
            route_name: 'finance.payments.list',
            active: false,
          },
        ]
      },
      {
        title: 'Observers',
        icon: 'fal fa-heart-rate',
        active: false,
        items: [
          {
            title: 'Listeners',
            icon: 'fal fa-function',
            route_name: 'listeners.list',
            active: false,
          },
          {
            title: 'Events',
            icon: 'fal fa-calendar-week',
            route_name: 'events.list',
            active: false,
            system: true,
          },
          {
            title: 'Handlers',
            icon: 'fal fa-lambda',
            route_name: 'handlers.list',
            active: false,
            system: true,
          },
        ]
      },
      {
        title: 'Интеграции',
        icon: 'fal fa-cogs',
        route_name: 'integrations',
        active: false,
        items: [
          {
            title: 'REST API',
            icon: 'fal fa-chart-network',
            route_name: 'integrations.api',
            active: false,
          },
          {
            title: 'WebHook',
            icon: 'fal fa-chart-network',
            route_name: 'integrations.webHooks',
            active: false,
          },
          {
            title: 'Методы доставки',
            icon: 'fal fa-pallet',
            route_name: 'delivery.methods',
            active: false,
          },
          {
            title: 'Методы оплаты',
            icon: 'fal fa-cash-register',
            route_name: 'finance.payment-method.list',
            active: false,
          },
        ]
      },
    ],

    miniVariant: false,
    expandOnHover: false,
  }),
  mounted() {
  },
  methods: {
    isDisabledMenuItem: function (obj) {
      return obj.hasOwnProperty('disabled') && obj['disabled']
    }
  },
  computed: {
    bg() {
      return this.background ? 'https://cdn.vuetifyjs.com/images/backgrounds/bg-2.jpg' : undefined
    },

  },

}
</script>

<style lang="scss">
.v-navigation-drawer .v-list-item {
  padding-left: 14px !important;
  padding-right: 14px !important;
}

.nav-active-item {
  /*color: #95ffff !important;*/
  color: white !important;
  background-color: rgba(#cae8e8, .1);

  /*background-color: transparent !important;*/
  font-weight: 600 !important;
}

.v-list-group--active {
  border-bottom: 1px solid #383c3e !important; // нижнее подчеркивание при открытом меню <v-list-group>
}

.exact-nav-list-item--active {
  color: #95ffff !important;
  /*color:white;*/
  background-color: transparent !important;

  &::before {
    background-color: transparent !important;
  }
}

.v-list-item--active.v-list-item.v-list-item--link.theme--dark {
  color: white;
}

.theme--dark.v-list-item--active:hover::before, .theme--dark.v-list-item--active::before {
  /*opacity: 0;*/
}

.navigationDrawerClass {
  /*background-image: linear-gradient(to right, #051d39, #1e344d);*/
  /*background-color: #13202f;*/
  /*background-image: linear-gradient(to right, #13202f, #203752);*/
}

</style>