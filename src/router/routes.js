
const routes = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/IndexPage.vue') },
      { path: '/kurs/:kurs_id', component: () => import('pages/KursOverviewPage.vue') },
      { path: '/do_kurs/:user_name/:kurs_id', component: () => import('pages/DoKursPage.vue') },
      { path: '/diagram', component: () => import('pages/DiagramPage.vue') },
      { path: '/statistics', component: () => import('pages/StatisticsPage.vue') },
      { path: '/admin', component: () => import('pages/AdminPage.vue') },

    ]
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue')
  }
]

export default routes
