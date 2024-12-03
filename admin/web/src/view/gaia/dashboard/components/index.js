import GvaCard from "./card.vue"
import GvaChart from "./charts.vue"
import GaiaAppQuotaTable from "./appTokenQuota.vue"
import GaiaAccountMoneyTable from "./accountMoneyTable.vue"
import GaiaAppTokenQuotaTable from "./appTokenQuotaTable.vue"

export {
    GvaCard,
    GvaChart,
    GaiaAccountMoneyTable,
    GaiaAppTokenQuotaTable,
    GaiaAppQuotaTable,
}

// 保留一位小数（不四舍五入）
export const truncateToOneDecimal = (value) => {
    if (typeof value !== 'number') {
        value = Number(value);
    }
    return Math.floor(value * 10) / 10;
}
// 根据使用率显示不同颜色
export const getColorClass = (used, limit) => {
    if (limit === -1) return ''; // 无限制时不应用颜色
    const ratio = used / limit;
    if (ratio === 1) return 'deep-red';
    if (ratio > 0.9) return 'light-red';
    if (ratio > 0.8) return 'yellow';
    return 'green';
}


// 根据mode值显示不同应用类型文本
export const getAppModeText = (mode) => {
    switch (mode){
        case 'chat':
            return '聊天助手';
        case 'advanced-chat':
            return '聊天助手(工作流)';
        case 'workflow':
            return '工作流';
        case 'agent-chat':
            return 'Agent';
        case 'completion':
            return '文本生成';
    }
    return '未知';
}

export const getAppModeColor = (mode) => {
    switch (mode) {
        case 'chat':
            return '#4CAF50'; // 绿色
        case 'advanced-chat':
            return '#f321b8'; // 粉色
        case 'workflow':
            return '#FFC107'; // 琥珀色
        case 'agent-chat':
            return '#9C27B0'; // 紫色
        case 'completion':
            return '#FF5722'; // 深橙色
        default:
            return '#757575'; // 灰色，用于未知类型
    }
}

// 根据mode值显示不同作图功能名
export const getAiImagePathText = (mode) => {
    switch (mode){
        case 'ai-draw':
            return 'AI抠图';
        case 'linkfox':
            return 'LinkFox商品合成';
        case 'ai-platform':
            return '以图搜图';
    }
    return '未知';
}

export const getAiImagePathColor = (mode) => {
    switch (mode) {
        case 'ai-draw':
            return '#4CAF50'; // 绿色
        case 'linkfox':
            return '#f321b8'; // 粉色
        case 'ai-platform':
            return '#FFC107'; // 琥珀色
        default:
            return '#757575'; // 灰色，用于未知类型
    }
}