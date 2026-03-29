# TopViewGame - Master Project Document

> **Last Updated:** 2025-10-18  
> **Current Phase:** Architecture Refactoring & Enhancement  
> **Project Status:** 🟢 Healthy - Functional game with solid improvement roadmap

---

## 📋 Executive Summary

**TopViewGame** is a functional roguelike game built with Go + Ebiten, featuring ECS architecture, procedural generation, and turn-based combat. The project has completed most foundational improvements and is ready for architectural refactoring.

**Key Achievements:**
- ✅ Functional ECS-based game engine
- ✅ Enhanced auto-movement system with visual feedback
- ✅ 87.5% MVP completion (7/8 items)
- ✅ Comprehensive improvement roadmap

**Next Priority:** Complete unit testing, begin modularization

---

## 🎯 Current Status & Metrics

### **Functional Features**
- ✅ **Core Engine:** Ebiten v2.8.8 + ECS (bytearena/ecs v1.0.0)
- ✅ **Gameplay:** Turn-based combat, A* pathfinding, FOV system
- ✅ **Generation:** Procedural dungeons with rooms/corridors
- ✅ **AI:** Monster pathfinding (Skeleton, Zombie types)
- ✅ **Visuals:** Tile-based rendering, animations, shaders
- ✅ **UX:** Auto-movement with 200ms delays, path visualization

### **Technical Specifications**
- **Language:** Go 1.23.8
- **Engine:** Ebiten v2.8.8 (2D game engine)
- **Architecture:** ECS with manual system orchestration
- **Map Size:** 80x50 tiles, 32x32 pixel tiles
- **Screen:** 40x25 visible tiles + UI
- **Dependencies:** norendren/go-fov, golang.org/x/image

---

## 🏗️ Architecture Analysis

### **Current Structure (After World Facade Refactor)**
- **Core Package:** `internal/core` - Position, Rect, GameData types
- **World Package:** `internal/world` - ECS facade with clean API  
- **Main Package:** Game logic with reduced ECS coupling
- **Facade Pattern:** World interface hides ECS complexity

### **World Facade Pattern Implementation**
```go
// Before: Direct ECS exposure
g.WorldManager.Manager.Query(g.WorldManager.Tags["players"])
pos := result.Components[g.WorldManager.PositionC].(*Position)

// After: Clean facade API  
for _, player := range g.World.QueryPlayers() {
    pos := g.World.GetPosition(player)
}
```

**Benefits Achieved:**
- **Encapsulation:** ECS manager hidden behind World interface
- **Type Safety:** Typed methods like `QueryPlayers()`, `GetHealth()`
- **Maintainability:** Easy to swap ECS libraries without client changes
- **Clean API:** `g.World.QueryMonsters()` vs `g.WorldManager.Manager.Query(tags)`

### **Current Structure & Responsibilities**

#### **Package Organization:**
- **`main` (root):** Core game types (Game, GameData, GameMap, Level), ECS components, and all systems (system_*.go). Acts as both application root and domain layer.
- **`controller/`:** Input handling abstractions (currently only human controller)
- **`event/`:** Lightweight input event type definitions
- **`state/`:** Animation/state machine primitives for player and monsters
- **`rand/`:** Custom linear congruential RNG with buffering

#### **Runtime Flow:**
1. `main()` constructs Game, initializes GameMap, Level, ECS world, and controller
2. Ebiten game loop invokes Update/Draw with manual system ordering: animations → cursor → auto-move → player/monster updates
3. ECS manages entities/components with tags for query pre-filtering

#### **Key Dependencies & Coupling Points:**
- **Tight Game coupling:** Game exposes many fields (Map, World, WorldTags, AutoMove*) accessed directly by systems
- **Dual spatial ownership:** GameMap.monsterPositions manually synchronized with ECS positions
- **Global component pointers:** positionC, renderableC, etc. declared in world.go, coupling initialization order
- **Asset loading:** Rendering systems depend on Ebiten assets and global helpers

---

## ⚠️ Issues & Technical Debt

### **Architectural Issues**
1. **Monolithic main package** - All domain types and systems co-located (36 Go files)
2. **Global component variables** - positionC, monsterC, etc. defined globally, hindering testing
3. **Manual monster position sync** - GameMap.monsterPositions must mirror ECS data
4. **~~Inconsistent coordinate handling~~** - ✅ **FIXED** in MVP
5. **Asset loading during gameplay** - Controllers/PathVisualizer load images on first use
6. **Pathfinding limitations** - AStar ignores occupancy, only checks static WALL tiles
7. **Non-reproducible RNG** - Time-based seeding, huge buffer (~64MB)
8. **~~UI text accumulation~~** - ✅ **FIXED** in MVP
9. **Missing test coverage** - No unit/integration tests
10. **Error handling** - Some log.Fatal calls remain, though improved
11. **Magic numbers** - Tile sizes, timing scattered throughout code
12. **Frame-dependent turns** - TurnCounter increments every frame, not time-based
13. **State package coupling** - Circular references between animation states and owners

### **Gameplay Limitations**
- Simple melee-only combat system
- No inventory/item system
- Limited monster types (2 variants)
- Single-level dungeons
- No character progression system

---

## 🚀 Improvement Programs

### **Minimum Viable Program (Quick Wins)**

| Item | Status | Description |
|------|--------|-------------|
| 1. Fix Level.InBounds | ✅ **DONE** | Upper bounds corrected to use >= |
| 2. Replace hardcoded tile size | ✅ **DONE** | DrawRenderables uses g.gd.TileWidth/Height |
| 3. ECS component initialization | ✅ **DONE** | resetComponentGlobals() prevents re-init |
| 4. Monster position sync | ✅ **DONE** | updateMonsterPosition() helper implemented |
| 5. AStar bounds checks | ✅ **DONE** | Neighbor expansion has proper bounds checking |
| 6. **Unit tests for AStar/TileLine** | ❌ **PENDING** | Add _test.go files for pathfinding |
| 7. Asset loading error handling | ✅ **DONE** | PathVisualizer errors handled gracefully |
| 8. Clamp lastText length | ✅ **DONE** | User log limited to 5 messages |

**MVP Completion: 87.5% (7/8 items)**

### **Maximum Improvement Program (Long-term Roadmap)**

#### **Phase 1: Foundation (Months 1-2)**
1. **Modularize packages** 🎯 **IN PROGRESS**
   - **1a.** ✅ **DONE** Extract data types (Position, Rect, GameData) → `internal/core`
   - **1b.** ✅ **DONE** Move ECS initialization → `internal/world` package
   - **1c.** Shift rendering systems → `internal/render` (incremental: DrawRenderables → HUD → UserLog → Cursor/Path)
   - **1d.** Relocate controller implementation with interface injection

2. **Encapsulate ECS world state**
   - Create WorldContext struct with manager, tags, component definitions
   - Eliminate global component variables

3. **Game loop scheduler**
   - Replace manual system ordering with configurable registry
   - Enable easier testing and system reordering

#### **Phase 2: Core Systems (Months 3-4)**
4. **Deterministic RNG service**
   - Injectable RNG interface for reproducible runs/tests
   - Reduce buffer size, move to dedicated package

5. **Pathfinding & occupancy improvements**
   - **5a.** Extend AStar with occupancy callback for dynamic obstacles
   - **5b.** Shared path visualization module (hover previews + auto-move)

6. **Saveable game state**
   - Design serialization layer (player stats, map layout)
   - Decouple Ebiten images from game logic

#### **Phase 3: Quality & Testing (Months 5-6)**
7. **Automated test suite**
   - Unit tests: map generation, FOV, combat, pathfinding
   - Integration tests: turn transitions, ECS operations
   - CI workflow setup

8. **Animation/state machine refactor**
   - Replace stater interface with explicit animation data structures
   - Consider data-driven animation system

#### **Phase 4: Advanced Features (Months 7+)**
9. **Controller abstraction**
   - Interface for input devices (AI, network controllers)
   - Decouple drawing from input detection

10. **Asset/resource management**
    - Centralized loading with caching and error reporting
    - Pre-load at startup with fallbacks

11. **Configuration & dependency injection**
    - GameConfig struct for dimensions, tile sizes
    - Enable easier scaling/resolution adjustments

12. **Performance profiling & optimization**
    - Frame time instrumentation
    - Chunked rendering, pooled slices for large maps

13. **Documentation & tooling**
    - Architecture documentation, UML diagrams
    - Code generation for ECS components

---

## 📊 Development Log & Milestones

### **Recent Achievements (2025-10)**
- ✅ **Enhanced auto-movement system** - 200ms delays, visual path feedback
- ✅ **Architecture analysis completed** - Comprehensive technical debt assessment
- ✅ **MVP program 87.5% complete** - 7/8 foundational improvements done
- ✅ **Improvement roadmap finalized** - Detailed phase-based plan
- ✅ **ECS modularization Phase 1a-1b** - Extracted core types and world management
- ✅ **World Facade Pattern implemented** - Clean API over ECS complexity

### **Current Sprint Goals**
- 🎯 **Complete MVP Item #6** - Add unit tests for AStar.GetPath and TileLine
- 🎯 **Continue Phase 1.1c** - Extract rendering systems to internal/render package
- 🎯 **Phase 1.1d** - Relocate controller implementation with interface injection

### **Key Metrics**
- **Codebase:** 36 Go files, ~15K lines
- **Architecture debt:** Well-identified and planned
- **Test coverage:** 0% → Target: 70%+ after Phase 3
- **Package coupling:** High → Target: Modular after Phase 1

---

## 🎮 Gameplay Vision & Future Features

### **Immediate Gameplay Enhancements**
- **Inventory system** - Item pickup, equipment, consumables
- **Extended combat** - Ranged weapons, magic spells
- **Monster variety** - More types with unique behaviors
- **Character progression** - Leveling, skills, attributes

### **Long-term Vision**
- **Multi-level dungeons** - Stairs, level transitions
- **Procedural content** - Varied biomes, special rooms
- **Save/load system** - Persistent game state
- **Modding support** - Data-driven content system

---

## 🔧 Development Guidelines

### **Code Standards**
- **ECS Patterns:** Components as data, systems as behavior
- **Error Handling:** Graceful degradation, avoid log.Fatal in runtime
- **Testing:** Unit tests for all core algorithms
- **Documentation:** Inline comments for complex algorithms

### **Architecture Principles**
- **Dependency Inversion:** Inject interfaces, not concrete types
- **Single Responsibility:** Each system handles one concern
- **Incremental Refactoring:** Maintain compilable state between changes
- **Performance Awareness:** Profile before optimizing

### **Git Workflow**
- **Feature branches** for each improvement program item
- **Atomic commits** with clear descriptions
- **CI/CD** integration after test suite implementation

---

## 📚 Resources & References

- **Ebiten Documentation:** https://ebiten.org/
- **ECS Pattern:** https://github.com/bytearena/ecs
- **Go FOV Library:** https://github.com/norendren/go-fov
- **Roguelike Development:** https://www.fatoldyeti.com/

---

## 🤖 AI Agent Memory Tags

`#roguelike #go #ebiten #ecs #architecture-refactoring #mvp-87-complete #pathfinding #auto-movement #modularization-ready #test-coverage-needed #phase1-next`

**Context for AI:** This project has solid foundations with recent auto-movement improvements. MVP is nearly complete (only unit tests remain). Ready to begin architectural modularization (Phase 1). All technical debt is well-documented and prioritized. Focus on incremental, compilable refactoring steps.
