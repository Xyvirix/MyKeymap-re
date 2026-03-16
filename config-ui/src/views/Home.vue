<script lang="ts" setup>
import { computed } from "vue"
import { useAppStore } from "@/store/app"
import { useConfigStore } from "@/store/config"

const appStore = useAppStore()
const configStore = useConfigStore()

const enabledCount = computed(() => configStore.enabledKeymaps.length)
const customCount = computed(() => configStore.customKeymaps.length)
const windowGroupCount = computed(() => configStore.options.windowGroups.length)
</script>

<template>
  <div class="home-grid">
    <section class="hero-panel">
      <div>
        <p class="hero-eyebrow">Preserve the power, modernize the desktop</p>
        <h2>让配置、驻留、热重载成为一个统一的桌面工作流</h2>
        <p>
          这个桌面版本保留原有映射逻辑和配置格式，只增强后台驻留、进程管理、保存/应用反馈和整体 UI 体验。
        </p>
      </div>

      <div class="hero-actions">
        <v-btn color="primary" prepend-icon="mdi-lightning-bolt" :disabled="!appStore.canApply" @click="configStore.applyConfig">
          应用当前修改
        </v-btn>
        <v-btn variant="outlined" color="secondary" prepend-icon="mdi-refresh" @click="appStore.restartEngine">
          重启引擎
        </v-btn>
      </div>
    </section>

    <section class="metric-grid">
      <article class="metric-card">
        <span>启用模式</span>
        <strong>{{ enabledCount }}</strong>
      </article>
      <article class="metric-card">
        <span>自定义模式</span>
        <strong>{{ customCount }}</strong>
      </article>
      <article class="metric-card">
        <span>窗口组</span>
        <strong>{{ windowGroupCount }}</strong>
      </article>
      <article class="metric-card">
        <span>引擎状态</span>
        <strong>{{ appStore.engineLabel }}</strong>
      </article>
    </section>

    <section class="guide-panel">
      <div class="panel-heading">
        <div>
          <p class="panel-eyebrow">Guide</p>
          <h3>使用说明与配置文档</h3>
        </div>
      </div>
      <div class="guide-frame">
        <iframe src="/config_doc.html" frameborder="0" width="100%" height="100%"></iframe>
      </div>
    </section>
  </div>
</template>

<style scoped>
.home-grid {
  display: grid;
  gap: 14px;
  min-width: 0;
}

.hero-panel,
.guide-panel {
  padding: 18px;
  border: 1px solid rgba(131, 149, 186, 0.18);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.86);
  box-shadow: 0 12px 30px rgba(44, 72, 122, 0.08);
}

.hero-panel {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  background:
    radial-gradient(circle at top right, rgba(76, 111, 255, 0.18), transparent 28%),
    rgba(255, 255, 255, 0.84);
}

.hero-eyebrow,
.panel-eyebrow {
  margin: 0 0 8px;
  color: #5c6b84;
  font-size: 0.8rem;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.hero-panel h2,
.panel-heading h3 {
  margin: 0;
  color: #172235;
}

.hero-panel p:last-of-type {
  max-width: 760px;
  margin: 10px 0 0;
  color: #56657f;
  line-height: 1.58;
}

.hero-actions {
  display: flex;
  flex-shrink: 0;
  gap: 10px;
}

.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.metric-card {
  padding: 16px 18px;
  border-radius: 18px;
  background: linear-gradient(180deg, #ffffff 0%, #f4f7fb 100%);
  border: 1px solid rgba(131, 149, 186, 0.16);
}

.metric-card span {
  display: block;
  margin-bottom: 8px;
  color: #7b899d;
  font-size: 0.84rem;
}

.metric-card strong {
  color: #18253a;
  font-size: 1.36rem;
}

.guide-frame {
  height: min(48vh, 620px);
  min-height: 360px;
  overflow: hidden;
  border-radius: 16px;
}

.guide-frame iframe {
  display: block;
}

.panel-heading {
  margin-bottom: 10px;
}

@media (max-width: 1320px) {
  .hero-panel {
    flex-direction: column;
    align-items: flex-start;
  }

  .metric-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
