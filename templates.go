package main

// ==================== ДАШБОРД ====================
const dashboardHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>СтроМенеджер — Дашборд</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Space+Grotesk:wght@500;600;700&display=swap');
        body { font-family: 'Inter', sans-serif; }
        .title { font-family: 'Space Grotesk', sans-serif; }
    </style>
</head>
<body class="bg-slate-950 text-slate-100 min-h-screen">
    <div class="max-w-7xl mx-auto p-8">
        <div class="flex justify-between items-center mb-12">
            <div>
                <h1 class="text-5xl font-bold title tracking-tight">СтроМенеджер</h1>
                <p class="text-slate-400 text-lg">Управление строительством и проектированием</p>
            </div>
            <a href="/projects" 
               class="bg-orange-500 hover:bg-orange-600 px-8 py-4 rounded-2xl font-medium flex items-center gap-3 transition-all">
                <i class="fas fa-list"></i> Все проекты
            </a>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400 text-sm">Всего проектов</div>
                <div class="text-6xl font-bold mt-4">{{.Total}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400 text-sm">Активных проектов</div>
                <div class="text-6xl font-bold mt-4 text-orange-400">{{.Active}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400 text-sm">Общий бюджет</div>
                <div class="text-6xl font-bold mt-4">{{.Budget}} ₽</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400 text-sm">Потрачено всего</div>
                <div class="text-6xl font-bold mt-4 text-emerald-400">{{.Spent}} ₽</div>
                <div class="mt-6 h-3 bg-slate-700 rounded-full overflow-hidden">
                    <div class="h-full bg-emerald-500 rounded-full transition-all" 
                         style="width: {{.SpentPercent}}%"></div>
                </div>
                <div class="text-xs text-slate-400 mt-2">{{.SpentPercent}}% от общего бюджета</div>
            </div>
        </div>

        <div class="text-center text-2xl font-medium text-slate-300">
            Средний прогресс по всем проектам: 
            <span class="text-orange-400 font-semibold">{{.AvgProgress}}%</span>
        </div>
    </div>
</body>
</html>`

// ==================== СПИСОК ПРОЕКТОВ ====================
const projectsHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Проекты — СтроМенеджер</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Space+Grotesk:wght@500;600;700&display=swap');
        body { font-family: 'Inter', sans-serif; }
        .title { font-family: 'Space Grotesk', sans-serif; }
    </style>
</head>
<body class="bg-slate-950 text-slate-100">
    <div class="max-w-7xl mx-auto p-8">
        <a href="/" class="inline-flex items-center gap-2 text-orange-400 hover:text-orange-300 mb-8">
            ← На дашборд
        </a>

        <div class="flex justify-between items-center mb-10">
            <h1 class="text-4xl font-semibold title">Все проекты</h1>
            <button onclick="document.getElementById('addProjectModal').classList.toggle('hidden')" 
                    class="bg-orange-500 hover:bg-orange-600 px-6 py-3 rounded-2xl flex items-center gap-2 font-medium">
                <i class="fas fa-plus"></i> Новый проект
            </button>
        </div>

        <div class="grid gap-6">
            {{range .}}
            <div class="bg-slate-900 rounded-3xl p-7 flex flex-col md:flex-row gap-6 items-center">
                <div class="flex-1">
                    <h3 class="text-2xl font-medium">{{.Name}}</h3>
                    <p class="text-slate-400 text-sm mt-1 line-clamp-2">{{.Description}}</p>
                    <div class="mt-6 h-3 bg-slate-700 rounded-full overflow-hidden">
                        <div class="h-3 bg-orange-500 rounded-full transition-all" 
                             style="width: {{.Progress}}%"></div>
                    </div>
                </div>
                <div class="text-right min-w-[120px]">
                    <div class="text-5xl font-bold text-orange-400">{{.Progress}}%</div>
                    <div class="text-xs text-slate-400">прогресс</div>
                    <div class="mt-2 text-lg font-medium">{{.Budget}} ₽</div>
                </div>
                <a href="/project/{{.ID}}" 
                   class="bg-orange-500 hover:bg-orange-600 px-8 py-4 rounded-2xl font-medium whitespace-nowrap">
                    Открыть проект
                </a>
            </div>
            {{end}}
        </div>

        <!-- Модальное окно добавления проекта -->
        <div id="addProjectModal" class="hidden fixed inset-0 bg-black/80 flex items-center justify-center z-50">
            <div class="bg-slate-900 rounded-3xl p-8 w-full max-w-lg mx-4">
                <h2 class="text-2xl font-semibold mb-6">Создать новый проект</h2>
                <form method="POST" action="/projects" class="space-y-5">
                    <input type="text" name="name" placeholder="Название проекта" required
                           class="w-full bg-slate-800 border border-slate-700 rounded-2xl px-5 py-4 focus:outline-none focus:border-orange-500">
                    <textarea name="description" placeholder="Краткое описание" rows="3"
                              class="w-full bg-slate-800 border border-slate-700 rounded-2xl px-5 py-4 focus:outline-none focus:border-orange-500"></textarea>
                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label class="text-xs text-slate-400 block mb-1">Дата начала</label>
                            <input type="date" name="start_date" class="w-full bg-slate-800 border border-slate-700 rounded-2xl px-5 py-4">
                        </div>
                        <div>
                            <label class="text-xs text-slate-400 block mb-1">Дата окончания</label>
                            <input type="date" name="end_date" class="w-full bg-slate-800 border border-slate-700 rounded-2xl px-5 py-4">
                        </div>
                    </div>
                    <input type="number" name="budget" placeholder="Бюджет проекта (в рублях)" step="10000"
                           class="w-full bg-slate-800 border border-slate-700 rounded-2xl px-5 py-4">
                    <div class="flex gap-4 pt-4">
                        <button type="submit" 
                                class="flex-1 bg-orange-500 hover:bg-orange-600 py-4 rounded-2xl font-medium transition-colors">
                            Создать проект
                        </button>
                        <button type="button" onclick="document.getElementById('addProjectModal').classList.add('hidden')"
                                class="flex-1 bg-slate-700 hover:bg-slate-600 py-4 rounded-2xl font-medium">
                            Отмена
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</body>
</html>`

// ==================== ДЕТАЛЬНАЯ СТРАНИЦА ПРОЕКТА ====================
const projectDetailHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Project.Name}} — СтроМенеджер</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
    <style>
        @import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=Space+Grotesk:wght@500;600;700&display=swap');
        body { font-family: 'Inter', sans-serif; }
        .title { font-family: 'Space Grotesk', sans-serif; }
    </style>
</head>
<body class="bg-slate-950 text-slate-100">
    <div class="max-w-7xl mx-auto p-8">
        <a href="/projects" class="inline-flex items-center gap-2 text-orange-400 hover:text-orange-300 mb-6">
            ← Все проекты
        </a>

        <div class="flex justify-between items-start">
            <div>
                <h1 class="text-4xl font-bold title">{{.Project.Name}}</h1>
                <p class="text-slate-400 mt-2 max-w-2xl">{{.Project.Description}}</p>
            </div>
            <div class="text-right">
                <div class="text-5xl font-bold text-orange-400">{{.Project.Progress}}%</div>
                <div class="text-sm text-slate-400">общий прогресс</div>
            </div>
        </div>

        <!-- Основные показатели проекта -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-10">
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="flex justify-between">
                    <div>
                        <div class="text-slate-400">Бюджет</div>
                        <div class="text-4xl font-semibold mt-2">{{.Project.Budget}} ₽</div>
                    </div>
                    <i class="fas fa-ruble-sign text-4xl text-slate-600"></i>
                </div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="flex justify-between">
                    <div>
                        <div class="text-slate-400">Потрачено</div>
                        <div class="text-4xl font-semibold mt-2 text-emerald-400">{{.Project.Spent}} ₽</div>
                    </div>
                    <i class="fas fa-coins text-4xl text-emerald-600"></i>
                </div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="flex justify-between">
                    <div>
                        <div class="text-slate-400">Сроки</div>
                        <div class="text-2xl font-medium mt-2">{{.Project.StartDate}} — {{.Project.EndDate}}</div>
                    </div>
                    <i class="fas fa-calendar-alt text-4xl text-slate-600"></i>
                </div>
            </div>
        </div>

        <!-- Объекты строительства -->
        <div class="mt-16">
            <div class="flex justify-between items-center mb-6">
                <h2 class="text-2xl font-semibold">Объекты строительства</h2>
                <button onclick="document.getElementById('addObjectModal').classList.toggle('hidden')" 
                        class="bg-orange-500 hover:bg-orange-600 px-6 py-3 rounded-2xl text-sm font-medium flex items-center gap-2">
                    <i class="fas fa-plus"></i> Добавить объект
                </button>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {{range .Objects}}
                <div class="bg-slate-900 rounded-3xl p-6">
                    <div class="flex justify-between items-start">
                        <div>
                            <h3 class="font-semibold text-xl">{{.Name}}</h3>
                            <p class="text-slate-400 text-sm">{{.Type}}</p>
                        </div>
                        <span class="text-xs bg-slate-800 px-3 py-1 rounded-full">{{.Area}} м²</span>
                    </div>
                    <div class="mt-6 h-2.5 bg-slate-700 rounded-full overflow-hidden">
                        <div class="h-full bg-orange-500 rounded-full" style="width: {{.Progress}}%"></div>
                    </div>
                    <div class="flex justify-between text-sm mt-3">
                        <span>{{.Budget}} ₽</span>
                        <span class="text-emerald-400">{{.Spent}} ₽</span>
                    </div>
                    <form action="/object/delete/{{.ID}}" method="POST" class="mt-6">
                        <button type="submit" onclick="return confirm('Удалить объект?')" 
                                class="text-red-500 text-sm hover:text-red-400">Удалить объект</button>
                    </form>
                </div>
                {{end}}
            </div>
        </div>

        <!-- График работ / Задачи -->
        <div class="mt-16">
            <div class="flex justify-between items-center mb-6">
                <h2 class="text-2xl font-semibold">График работ</h2>
                <button onclick="document.getElementById('addTaskModal').classList.toggle('hidden')" 
                        class="bg-orange-500 hover:bg-orange-600 px-6 py-3 rounded-2xl text-sm font-medium flex items-center gap-2">
                    <i class="fas fa-plus"></i> Добавить работу
                </button>
            </div>

            <div class="space-y-5">
                {{range .Tasks}}
                <div class="bg-slate-900 rounded-3xl p-6 flex items-center gap-6">
                    <div class="flex-1">
                        <div class="font-medium text-lg">{{.Name}}</div>
                        <div class="text-slate-400 text-sm mt-1">
                            {{.AssignedTo}} • {{.StartDate}} — {{.EndDate}}
                        </div>
                        <div class="mt-4 h-2 bg-slate-700 rounded-full overflow-hidden">
                            <div class="h-2 bg-blue-500 rounded-full" style="width: {{.Progress}}%"></div>
                        </div>
                    </div>
                    <div class="text-right min-w-[140px]">
                        <div class="text-3xl font-bold text-blue-400">{{.Progress}}%</div>
                        <div class="text-xs text-slate-400 mt-1">
                            {{.Spent}} / {{.Estimated}} ₽
                        </div>
                    </div>
                    <form action="/task/delete/{{.ID}}" method="POST">
                        <button type="submit" onclick="return confirm('Удалить работу?')" 
                                class="text-red-500 hover:text-red-400 px-4 py-2">✕</button>
                    </form>
                </div>
                {{end}}
            </div>
        </div>
    </div>

    <!-- Модальное окно добавления объекта -->
    <div id="addObjectModal" class="hidden fixed inset-0 bg-black/80 flex items-center justify-center z-50">
        <div class="bg-slate-900 rounded-3xl p-8 w-full max-w-md mx-4">
            <h2 class="text-2xl font-semibold mb-6">Новый объект</h2>
            <form method="POST" action="/add-object" class="space-y-5">
                <input type="hidden" name="project_id" value="{{.Project.ID}}">
                <input type="text" name="name" placeholder="Название объекта" required
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4">
                <input type="text" name="type" placeholder="Тип (здание, фундамент, парковка и т.д.)"
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4">
                <input type="number" name="area" placeholder="Площадь (м²)" step="0.1"
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4">
                <input type="number" name="budget" placeholder="Бюджет объекта (₽)" step="10000"
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4">
                <div class="flex gap-4">
                    <button type="submit" class="flex-1 bg-orange-500 hover:bg-orange-600 py-4 rounded-2xl font-medium">Добавить объект</button>
                    <button type="button" onclick="document.getElementById('addObjectModal').classList.add('hidden')"
                            class="flex-1 bg-slate-700 hover:bg-slate-600 py-4 rounded-2xl">Отмена</button>
                </div>
            </form>
        </div>
    </div>

    <!-- Модальное окно добавления задачи -->
    <div id="addTaskModal" class="hidden fixed inset-0 bg-black/80 flex items-center justify-center z-50">
        <div class="bg-slate-900 rounded-3xl p-8 w-full max-w-md mx-4">
            <h2 class="text-2xl font-semibold mb-6">Новая работа</h2>
            <form method="POST" action="/add-task" class="space-y-5">
                <input type="hidden" name="project_id" value="{{.Project.ID}}">
                <input type="text" name="name" placeholder="Название работы" required
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4">
                <input type="text" name="assigned_to" placeholder="Исполнитель (ФИО)" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4">
                <div class="grid grid-cols-2 gap-4">
                    <input type="date" name="start_date" class="bg-slate-800 rounded-2xl px-5 py-4">
                    <input type="date" name="end_date" class="bg-slate-800 rounded-2xl px-5 py-4">
                </div>
                <input type="number" name="estimated" placeholder="Плановые затраты (₽)" step="1000"
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4">
                <div class="flex gap-4">
                    <button type="submit" class="flex-1 bg-orange-500 hover:bg-orange-600 py-4 rounded-2xl font-medium">Добавить работу</button>
                    <button type="button" onclick="document.getElementById('addTaskModal').classList.add('hidden')"
                            class="flex-1 bg-slate-700 hover:bg-slate-600 py-4 rounded-2xl">Отмена</button>
                </div>
            </form>
        </div>
    </div>
</body>
</html>`