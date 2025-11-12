import pygame

def readState(fileName):
    try:
        file = open(fileName, "r")
        line = file.readline()
        L = int(line)
        cells = []
        for i in range(L):
            row = file.readline().split()
            cells.append(row)
        file.close()
        return cells
    except ValueError:
        print(f"Ошибка чтения {fileName}")
        
    except FileNotFoundError:
        print(f'Файл {fileName} не найден.')
    return []

def drawState(cells):
    L = len(cells)
    if L == 0:
        return
    w = WEIGHT//L
    h = HEIGHT//L
    for i in range(L):
        for j in range(L):
            if len(cells[i]) < L:
                return
            if cells[i][j] == "1":
                x = w * i
                y = h * (j) #L - j - 1
                pygame.draw.rect(screen, COLOR, ( y, x, w, h))
    

WEIGHT = 720
HEIGHT = 720
FPS = 1  # Скорость анимации (состояний в секунду) - 1 кадр в секунду
COLOR = (0, 171, 131)

pygame.init()
screen = pygame.display.set_mode((WEIGHT, HEIGHT))
pygame.display.set_caption("GoL Viewer - Анимация")
clock = pygame.time.Clock()

# Находим количество состояний
def getMaxStates():
    maxState = 0
    while True:
        filename = f"states/state_{maxState}.txt"
        try:
            file = open(filename, "r")
            file.close()
            maxState += 1
        except FileNotFoundError:
            break
    return maxState

maxStates = getMaxStates()
print(f"Найдено {maxStates} состояний")

if maxStates == 0:
    print("Файлы состояний не найдены. Запустите сначала template.go!")
    pygame.quit()
    exit()

currentState = 0
running = True

while running:
    clock.tick(FPS)
    screen.fill((255, 255, 255))
    
    # Читаем и отображаем текущее состояние
    filename = f"states/state_{currentState}.txt"
    cells = readState(filename)
    drawState(cells)
    pygame.display.set_caption(f"GoL Viewer - Итерация {currentState}/{maxStates-1}")
    pygame.display.update()
    
    # Переключаемся на следующее состояние
    currentState = (currentState + 1) % maxStates
    
    # Проверяем выход
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False
        
pygame.quit()

