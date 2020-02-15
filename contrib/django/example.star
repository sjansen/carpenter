url(
    "django.contrib.admin.sites.AdminSite.index",
    path = {
        "prefix": [
            'admin',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/": "/admin/",
    },
)

url(
    "django.contrib.admin.sites.AdminSite.login",
    path = {
        "prefix": [
            'admin',
            'login',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/login/": "/admin/login/",
    },
)

url(
    "django.contrib.admin.sites.AdminSite.logout",
    path = {
        "prefix": [
            'admin',
            'logout',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/logout/": "/admin/logout/",
    },
)

url(
    "django.contrib.admin.sites.AdminSite.password_change",
    path = {
        "prefix": [
            'admin',
            'password_change',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/password_change/": "/admin/password_change/",
    },
)

url(
    "django.contrib.admin.sites.AdminSite.password_change_done",
    path = {
        "prefix": [
            'admin',
            'password_change',
            'done',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/password_change/done/": "/admin/password_change/done/",
    },
)

url(
    "django.contrib.admin.sites.AdminSite.i18n_javascript",
    path = {
        "prefix": [
            'admin',
            'jsi18n',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/jsi18n/": "/admin/jsi18n/",
    },
)

url(
    "django.contrib.contenttypes.views.shortcut",
    path = {
        "prefix": [
            'admin',
            'r',
            (r"[0-9]+", 'CONTENT_TYPE_ID'),
            (r".+", 'OBJECT_ID'),
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/r/42/foo/": "/admin/r/CONTENT_TYPE_ID/OBJECT_ID/",
        "/admin/r/6/bar/": "/admin/r/CONTENT_TYPE_ID/OBJECT_ID/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.changelist_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'group',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/group/": "/admin/auth/group/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.add_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'group',
            'add',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/group/add/": "/admin/auth/group/add/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.autocomplete_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'group',
            'autocomplete',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/group/autocomplete/": "/admin/auth/group/autocomplete/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.history_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'group',
            (r".+", 'OBJECT_ID'),
            'history',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/group/foo/history/": "/admin/auth/group/OBJECT_ID/history/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.delete_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'group',
            (r".+", 'OBJECT_ID'),
            'delete',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/group/foo/delete/": "/admin/auth/group/OBJECT_ID/delete/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.change_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'group',
            (r".+", 'OBJECT_ID'),
            'change',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/group/foo/change/": "/admin/auth/group/OBJECT_ID/change/",
    },
)

url(
    "django.contrib.auth.admin.UserAdmin.user_change_password",
    path = {
        "prefix": [
            'admin',
            'auth',
            'user',
            (r"[^/]+", 'ID'),
            'password',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/user/sjansen/password/": "/admin/auth/user/ID/password/",
        "/admin/auth/user/9/password/": "/admin/auth/user/ID/password/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.changelist_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'user',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/user/": "/admin/auth/user/",
    },
)

url(
    "django.contrib.auth.admin.UserAdmin.add_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'user',
            'add',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/user/add/": "/admin/auth/user/add/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.autocomplete_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'user',
            'autocomplete',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/user/autocomplete/": "/admin/auth/user/autocomplete/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.history_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'user',
            (r".+", 'OBJECT_ID'),
            'history',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/user/foo/history/": "/admin/auth/user/OBJECT_ID/history/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.delete_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'user',
            (r".+", 'OBJECT_ID'),
            'delete',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/user/foo/delete/": "/admin/auth/user/OBJECT_ID/delete/",
    },
)

url(
    "django.contrib.admin.options.ModelAdmin.change_view",
    path = {
        "prefix": [
            'admin',
            'auth',
            'user',
            (r".+", 'OBJECT_ID'),
            'change',
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/user/foo/change/": "/admin/auth/user/OBJECT_ID/change/",
        "/admin/auth/user/bar/change/": "/admin/auth/user/OBJECT_ID/change/",
    },
)

url(
    "django.contrib.admin.sites.AdminSite.app_index",
    path = {
        "prefix": [
            'admin',
            (r"auth", 'APP_LABEL'),
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/auth/": "/admin/APP_LABEL/",
    },
)

url(
    "django.views.defaults.page_not_found",
    path = {
        "prefix": [
            'admin',
            (r"[^/]", 'APP_LABEL', r"auth|jsi18n|login|logout|password_change"),
        ],
        "suffix": "/",
    },
    query = {
        "other": "X",
    },
    tests = {
        "/admin/roles/": "/admin/APP_LABEL/",
    },
)

