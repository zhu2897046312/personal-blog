package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// GenerateCacheKey 生成缓存key
// 参数：
// - conditions: 查询条件map
// - args: 其他参数（如page, pageSize等）
func GenerateCacheKey(conditions map[string]interface{}, args ...interface{}) string {
	// 处理conditions
	var keys []string
	for k, v := range conditions {
		keys = append(keys, fmt.Sprintf("%s:%v", k, v))
	}
	// 对key进行排序，确保相同的条件生成相同的key
	sort.Strings(keys)
	condStr := strings.Join(keys, "_")

	// 处理其他参数
	var argStrs []string
	for _, arg := range args {
		argStrs = append(argStrs, fmt.Sprintf("%v", arg))
	}
	argStr := strings.Join(argStrs, "_")

	// 组合完整的key字符串
	fullStr := fmt.Sprintf("%s_%s", condStr, argStr)

	// 如果key太长，使用MD5生成固定长度的key
	if len(fullStr) > 100 {
		return generateMD5(fullStr)
	}

	return fullStr
}

// GenerateListCacheKey 生成列表缓存key
func GenerateListCacheKey(prefix string, page, pageSize int, filters ...interface{}) string {
	parts := []string{prefix}
	
	// 添加分页信息
	parts = append(parts, fmt.Sprintf("page:%d", page))
	parts = append(parts, fmt.Sprintf("size:%d", pageSize))
	
	// 添加过滤条件
	for _, filter := range filters {
		// 将过滤条件转换为JSON字符串
		if filterJSON, err := json.Marshal(filter); err == nil {
			parts = append(parts, string(filterJSON))
		}
	}
	
	// 组合key
	key := strings.Join(parts, "_")
	
	// 如果key太长，使用MD5
	if len(key) > 100 {
		return generateMD5(key)
	}
	
	return key
}

// GenerateDetailCacheKey 生成详情缓存key
func GenerateDetailCacheKey(prefix string, id interface{}, args ...interface{}) string {
	parts := []string{prefix, fmt.Sprintf("id:%v", id)}
	
	// 添加其他参数
	for _, arg := range args {
		parts = append(parts, fmt.Sprintf("%v", arg))
	}
	
	return strings.Join(parts, "_")
}

// generateMD5 生成MD5哈希
func generateMD5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

// CombineCacheKeys 组合多个缓存key
func CombineCacheKeys(keys ...string) string {
	return strings.Join(keys, ":")
}

// ParseCacheKey 解析缓存key
func ParseCacheKey(key string) map[string]string {
	result := make(map[string]string)
	parts := strings.Split(key, "_")
	
	for _, part := range parts {
		if strings.Contains(part, ":") {
			kv := strings.SplitN(part, ":", 2)
			if len(kv) == 2 {
				result[kv[0]] = kv[1]
			}
		}
	}
	
	return result
}
