package utils

import (
	"strings"
	"unicode"
)

// KanaMap содержит маппинг катаканы в хирагану для базовых символов
var KanaMap = map[rune]rune{
	// Катакана -> Хирагана
	'ァ': 'ぁ', 'ア': 'あ', 'ィ': 'ぃ', 'イ': 'い', 'ゥ': 'ぅ',
	'ウ': 'う', 'ェ': 'ぇ', 'エ': 'え', 'ォ': 'ぉ', 'オ': 'お',
	'カ': 'か', 'ガ': 'が', 'キ': 'き', 'ギ': 'ぎ', 'ク': 'く',
	'グ': 'ぐ', 'ケ': 'け', 'ゲ': 'げ', 'コ': 'こ', 'ゴ': 'ご',
	'サ': 'さ', 'ザ': 'ざ', 'シ': 'し', 'ジ': 'じ', 'ス': 'す',
	'ズ': 'ず', 'セ': 'せ', 'ゼ': 'ぜ', 'ソ': 'そ', 'ゾ': 'ぞ',
	'タ': 'た', 'ダ': 'だ', 'チ': 'ち', 'ヂ': 'ぢ', 'ッ': 'っ',
	'ツ': 'つ', 'ヅ': 'づ', 'テ': 'て', 'デ': 'で', 'ト': 'と',
	'ド': 'ど', 'ナ': 'な', 'ニ': 'に', 'ヌ': 'ぬ', 'ネ': 'ね',
	'ノ': 'の', 'ハ': 'は', 'バ': 'ば', 'パ': 'ぱ', 'ヒ': 'ひ',
	'ビ': 'び', 'ピ': 'ぴ', 'フ': 'ふ', 'ブ': 'ぶ', 'プ': 'ぷ',
	'ヘ': 'へ', 'ベ': 'べ', 'ペ': 'ぺ', 'ホ': 'ほ', 'ボ': 'ぼ',
	'ポ': 'ぽ', 'マ': 'ま', 'ミ': 'み', 'ム': 'む', 'メ': 'め',
	'モ': 'も', 'ヤ': 'や', 'ャ': 'ゃ', 'ユ': 'ゆ', 'ュ': 'ゅ',
	'ヨ': 'よ', 'ョ': 'ょ', 'ラ': 'ら', 'リ': 'り', 'ル': 'る',
	'レ': 'れ', 'ロ': 'ろ', 'ワ': 'わ', 'ヮ': 'ゎ', 'ヰ': 'ゐ',
	'ヱ': 'ゑ', 'ヲ': 'を', 'ン': 'ん',
	// Дополнительные символы
	'ヴ': 'ゔ', 'ヵ': 'ゕ', 'ヶ': 'ゖ', 'ヽ': 'ゝ', 'ヾ': 'ゞ',
}

// ReverseKanaMap содержит обратный маппинг (хирагана -> катакана)
var ReverseKanaMap = make(map[rune]rune)

func init() {
	// Инициализируем обратный маппинг
	for k, v := range KanaMap {
		ReverseKanaMap[v] = k
	}
}

// IsKatakana проверяет, является ли руна катаканой
func IsKatakana(r rune) bool {
	// Основной диапазон катаканы
	if r >= 'ァ' && r <= 'ヾ' {
		return true
	}
	// Расширенная катакана
	if r >= 'ㇰ' && r <= 'ㇿ' {
		return true
	}
	return false
}

// IsHiragana проверяет, является ли руна хираганой
func IsHiragana(r rune) bool {
	return r >= 'ぁ' && r <= 'ゖ'
}

// ToHiragana преобразует катакану в хирагану
func ToHiragana(s string) string {
	var result strings.Builder
	for _, r := range s {
		if mapped, ok := KanaMap[r]; ok {
			result.WriteRune(mapped)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToKatakana преобразует хирагану в катакану
func ToKatakana(s string) string {
	var result strings.Builder
	for _, r := range s {
		if mapped, ok := ReverseKanaMap[r]; ok {
			result.WriteRune(mapped)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// NormalizeOnReading нормализует онъёми для поиска
// Возвращает оба варианта: оригинальный и преобразованный
func NormalizeOnReading(on string) []string {
	variants := []string{on}

	// Проверяем, содержит ли строка катакану или хирагану
	hasKatakana := false
	hasHiragana := false

	for _, r := range on {
		if IsKatakana(r) {
			hasKatakana = true
		}
		if IsHiragana(r) {
			hasHiragana = true
		}
	}

	// Если содержит только катакану, добавляем хирагану
	if hasKatakana && !hasHiragana {
		variants = append(variants, ToHiragana(on))
	}

	// Если содержит только хирагану, добавляем катакану
	if hasHiragana && !hasKatakana {
		variants = append(variants, ToKatakana(on))
	}

	// Удаляем дубликаты
	return removeDuplicates(variants)
}

// removeDuplicates удаляет дубликаты из среза строк
func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if item == "" {
			continue
		}
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// IsJapaneseChar проверяет, является ли символ японским (канзи, кана)
func IsJapaneseChar(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana)
}
