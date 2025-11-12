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
        return cells
    except ValueError:
        print(line)
        
    except FileNotFoundError:
        print('Файл state.txt не найден.')
    finally:
        file.close()
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
FPS = 15
COLOR = (0, 171, 131)

pygame.init()
screen = pygame.display.set_mode((WEIGHT, HEIGHT))
pygame.display.set_caption("GoL Viewer")
clock = pygame.time.Clock()

# pygame.display.update()
running = True
while running:
    clock.tick(FPS)
    screen.fill((255, 255, 255))
    cells = readState("state_after_1_day.txt")
    drawState(cells)
    pygame.display.update()
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False
        
                    
pygame.quit()

