class DinosaurGame {
    constructor() {
        this.dinosaur = document.getElementById('dinosaur');
        this.gameArea = document.querySelector('.game-area');
        this.scoreElement = document.getElementById('score');
        this.gameOverElement = document.getElementById('gameOver');
        
        this.isJumping = false;
        this.isGameOver = false;
        this.score = 0;
        this.gameSpeed = 3;
        this.cacti = [];
        
        this.init();
    }
    
    init() {
        this.bindEvents();
        this.startGame();
    }
    
    bindEvents() {
        document.addEventListener('keydown', (e) => {
            if (this.isGameOver) {
                if (e.code === 'Space') {
                    this.restartGame();
                }
                return;
            }
            
            if (e.code === 'ArrowUp' || e.code === 'Space') {
                e.preventDefault();
                this.jump();
            }
        });
    }
    
    jump() {
        if (this.isJumping) return;
        
        this.isJumping = true;
        this.dinosaur.classList.add('jumping');
        
        setTimeout(() => {
            this.isJumping = false;
            this.dinosaur.classList.remove('jumping');
        }, 600);
    }
    
    createCactus() {
        const cactus = document.createElement('div');
        cactus.className = 'cactus';
        this.gameArea.appendChild(cactus);
        this.cacti.push({
            element: cactus,
            x: 800
        });
    }
    
    updateCacti() {
        this.cacti.forEach((cactus, index) => {
            cactus.x -= this.gameSpeed;
            cactus.element.style.right = (800 - cactus.x) + 'px';
            
            // 检查碰撞
            if (this.checkCollision(cactus)) {
                this.gameOver();
                return;
            }
            
            // 移除屏幕外的仙人掌
            if (cactus.x < -20) {
                cactus.element.remove();
                this.cacti.splice(index, 1);
                this.score += 10;
                this.scoreElement.textContent = this.score;
                
                // 增加游戏速度
                if (this.score % 100 === 0) {
                    this.gameSpeed += 0.5;
                }
            }
        });
    }
    
    checkCollision(cactus) {
        const dinosaurRect = this.dinosaur.getBoundingClientRect();
        const cactusRect = cactus.element.getBoundingClientRect();
        
        return !(dinosaurRect.right < cactusRect.left + 10 ||
                dinosaurRect.left > cactusRect.right - 10 ||
                dinosaurRect.bottom < cactusRect.top + 10 ||
                dinosaurRect.top > cactusRect.bottom - 10);
    }
    
    gameOver() {
        this.isGameOver = true;
        this.gameOverElement.style.display = 'block';
    }
    
    restartGame() {
        this.isGameOver = false;
        this.score = 0;
        this.gameSpeed = 3;
        this.cacti = [];
        this.scoreElement.textContent = '0';
        this.gameOverElement.style.display = 'none';
        
        // 清除所有仙人掌
        document.querySelectorAll('.cactus').forEach(cactus => cactus.remove());
        
        this.startGame();
    }
    
    startGame() {
        this.gameLoop();
        this.cactusSpawner();
    }
    
    gameLoop() {
        if (!this.isGameOver) {
            this.updateCacti();
            requestAnimationFrame(() => this.gameLoop());
        }
    }
    
    cactusSpawner() {
        if (this.isGameOver) return;
        
        // 随机生成仙人掌
        if (Math.random() < 0.005 + (this.score * 0.00001)) {
            this.createCactus();
        }
        
        setTimeout(() => {
            this.cactusSpawner();
        }, 100);
    }
}

// 启动游戏
window.addEventListener('load', () => {
    new DinosaurGame();
});