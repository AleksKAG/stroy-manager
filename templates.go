package main

// ==================== ДАШБОРД ====================
const dashboardHTML = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>СтроМенеджер</title>
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
                <p class="text-slate-400 text-lg">Управление проектированием и строительством</p>
            </div>
            <a href="/projects" class="bg-orange-500 hover:bg-orange-600 px-8 py-4 rounded-2xl font-medium flex items-center gap-3">
                <i class="fas fa-list"></i> Все проекты
            </a>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Всего проектов</div>
                <div class="text-6xl font-bold mt-4">{{.Total}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Активных</div>
                <div class="text-6xl font-bold mt-4 text-orange-400">{{.Active}}</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Общий бюджет</div>
                <div class="text-6xl font-bold mt-4">{{.Budget}} ₽</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Потрачено</div>
                <div class="text-6xl font-bold mt-4 text-emerald-400">{{.Spent}} ₽</div>
            </div>
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
</head>
<body class="bg-slate-950 text-slate-100">
    <div class="max-w-7xl mx-auto p-8">
        <a href="/" class="text-orange-400 mb-8 inline-flex items-center gap-2">← На дашборд</a>
        <h1 class="text-4xl font-semibold mb-8">Все проекты</h1>
        <!-- Список проектов (можно расширить позже) -->
    </div>
</body>
</html>`

// ==================== ДЕТАЛЬНАЯ СТРАНИЦА ПРОЕКТА С РЕДАКТИРОВАНИЕМ ====================
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
        <a href="/projects" class="inline-flex items-center gap-2 text-orange-400 mb-8">← Все проекты</a>

        <h1 class="text-4xl font-bold title">{{.Project.Name}}</h1>
        <p class="text-slate-400 mt-1">{{.Project.Description}}</p>

        <!-- Основные показатели -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mt-10">
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Бюджет</div>
                <div class="text-4xl font-semibold mt-3">{{.Project.Budget}} ₽</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Потрачено</div>
                <div class="text-4xl font-semibold mt-3 text-emerald-400">{{.Project.Spent}} ₽</div>
            </div>
            <div class="bg-slate-900 rounded-3xl p-8">
                <div class="text-slate-400">Прогресс проекта</div>
                <div class="text-5xl font-bold text-orange-400 mt-3">{{.Project.Progress}}%</div>
            </div>
        </div>

        <!-- Объекты строительства -->
        <div class="mt-16">
            <div class="flex justify-between items-center mb-6">
                <h2 class="text-2xl font-semibold">Объекты строительства</h2>
                <button onclick="showAddObjectModal()" 
                        class="bg-orange-500 hover:bg-orange-600 px-6 py-3 rounded-2xl flex items-center gap-2 text-sm font-medium">
                    <i class="fas fa-plus"></i> Добавить объект
                </button>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {{range .Objects}}
                <div class="bg-slate-900 rounded-3xl p-6">
                    <div class="flex justify-between items-start">
                        <div>
                            <h3 class="text-xl font-semibold">{{.Name}}</h3>
                            <p class="text-slate-400">{{.Type}} • {{.Floors}} этажей</p>
                            <p class="text-sm text-slate-500">{{.Material}}</p>
                        </div>
                        <button onclick="showEditObjectModal({{.ID}}, '{{.Name}}', '{{.Type}}', {{.Area}}, {{.Budget}}, {{.Floors}}, '{{.Material}}', '{{.Status}}')" 
                                class="text-orange-400 hover:text-white">
                            <i class="fas fa-edit"></i>
                        </button>
                    </div>

                    <div class="mt-6 h-2.5 bg-slate-700 rounded-full overflow-hidden">
                        <div class="h-full bg-orange-500 rounded-full" style="width: {{.Progress}}%"></div>
                    </div>

                    <div class="flex justify-between mt-4 text-sm">
                        <span class="font-medium">{{.Budget}} ₽</span>
                        <span class="text-emerald-400">{{.Spent}} ₽</span>
                    </div>

                    <form action="/object/delete/{{.ID}}" method="POST" class="mt-6">
                        <button type="submit" onclick="return confirm('Удалить объект?')" 
                                class="text-red-500 hover:text-red-400 text-sm">Удалить</button>
                    </form>
                </div>
                {{end}}
            </div>
        </div>
    </div>

    <!-- ==================== МОДАЛЬНОЕ ОКНО ДОБАВЛЕНИЯ ==================== -->
    <div id="addObjectModal" class="hidden fixed inset-0 bg-black/70 flex items-center justify-center z-50">
        <div class="bg-slate-900 rounded-3xl p-8 w-full max-w-lg mx-4">
            <h2 class="text-2xl font-semibold mb-6">Новый объект</h2>
            <form method="POST" action="/add-object">
                <input type="hidden" name="project_id" value="{{.Project.ID}}">

                <input type="text" name="name" placeholder="Название объекта" required 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="text" name="type" placeholder="Тип (здание, фундамент, парковка...)" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="number" name="area" placeholder="Площадь (м²)" step="0.1" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="number" name="budget" placeholder="Бюджет объекта (₽)" step="10000" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="number" name="floors" placeholder="Количество этажей" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="text" name="material" placeholder="Материал (железобетон, кирпич и т.д.)" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-6">

                <div class="flex gap-4">
                    <button type="submit" class="flex-1 bg-orange-500 hover:bg-orange-600 py-4 rounded-2xl font-medium">Добавить объект</button>
                    <button type="button" onclick="hideAddObjectModal()" 
                            class="flex-1 bg-slate-700 hover:bg-slate-600 py-4 rounded-2xl">Отмена</button>
                </div>
            </form>
        </div>
    </div>

    <!-- ==================== МОДАЛЬНОЕ ОКНО РЕДАКТИРОВАНИЯ ==================== -->
    <div id="editObjectModal" class="hidden fixed inset-0 bg-black/70 flex items-center justify-center z-50">
        <div class="bg-slate-900 rounded-3xl p-8 w-full max-w-lg mx-4">
            <h2 class="text-2xl font-semibold mb-6">Редактировать объект</h2>
            <form method="POST" action="/object/edit">
                <input type="hidden" name="id" id="edit_id">
                <input type="hidden" name="project_id" value="{{.Project.ID}}">

                <input type="text" name="name" id="edit_name" placeholder="Название объекта" required 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="text" name="type" id="edit_type" placeholder="Тип объекта" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="number" name="area" id="edit_area" placeholder="Площадь (м²)" step="0.1" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="number" name="budget" id="edit_budget" placeholder="Бюджет (₽)" step="10000" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="number" name="floors" id="edit_floors" placeholder="Количество этажей" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-4">
                <input type="text" name="material" id="edit_material" placeholder="Материал" 
                       class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-6">

                <select name="status" id="edit_status" 
                        class="w-full bg-slate-800 rounded-2xl px-5 py-4 mb-6">
                    <option value="in_progress">В работе</option>
                    <option value="completed">Завершён</option>
                    <option value="on_hold">Приостановлен</option>
                </select>

                <div class="flex gap-4">
                    <button type="submit" class="flex-1 bg-orange-500 hover:bg-orange-600 py-4 rounded-2xl font-medium">Сохранить изменения</button>
                    <button type="button" onclick="hideEditObjectModal()" 
                            class="flex-1 bg-slate-700 hover:bg-slate-600 py-4 rounded-2xl">Отмена</button>
                </div>
            </form>
        </div>
    </div>

    <script>
        function showAddObjectModal() {
            document.getElementById('addObjectModal').classList.remove('hidden');
        }
        function hideAddObjectModal() {
            document.getElementById('addObjectModal').classList.add('hidden');
        }

        function showEditObjectModal(id, name, type, area, budget, floors, material, status) {
            document.getElementById('edit_id').value = id;
            document.getElementById('edit_name').value = name;
            document.getElementById('edit_type').value = type;
            document.getElementById('edit_area').value = area;
            document.getElementById('edit_budget').value = budget;
            document.getElementById('edit_floors').value = floors;
            document.getElementById('edit_material').value = material;
            if (status) document.getElementById('edit_status').value = status;

            document.getElementById('editObjectModal').classList.remove('hidden');
        }

        function hideEditObjectModal() {
            document.getElementById('editObjectModal').classList.add('hidden');
        }
    </script>
</body>
</html>`