//go:build analyzer_angular

package analyze

import (
	"github.com/wintbiit/semantic-release-go/utils"
	"slices"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/rs/zerolog/log"
	"github.com/wintbiit/semantic-release-go/types"
)

const (
	FieldFeatureTitle     = "🌟 Features"
	FieldFixTitle         = "🐞 Bug Fixes"
	FieldPerformanceTitle = "🚀 Performance Improvements"
	FieldBreakingTitle    = "💥 Breaking Changes"
	FieldRefactorTitle    = "⚙️ Code Refactoring"
	FieldStyleTitle       = "🎨 UI"
	FieldTestTitle        = "🧪 Tests"
	FieldDocsTitle        = "📚 Documentation"
	FieldChoreTitle       = "🧹 Chores"
	FieldBuildTitle       = "🏗️ Build System"
	FieldRevertTitle      = "⏪ Reverts"
	FieldDependencyTitle  = "📦 Dependencies"
	FieldSecurityTitle    = "🔒 Security"
	FieldOtherTitle       = "🔧 Others"
)

var typeMap = map[string]string{
	"feat":    FieldFeatureTitle,
	"feature": FieldFeatureTitle,

	"fix":  FieldFixTitle,
	"bug":  FieldFixTitle,
	"bugs": FieldFixTitle,

	"perf":         FieldPerformanceTitle,
	"improvement":  FieldPerformanceTitle,
	"enhancement":  FieldPerformanceTitle,
	"enhance":      FieldPerformanceTitle,
	"optimize":     FieldPerformanceTitle,
	"optimization": FieldPerformanceTitle,
	"performance":  FieldPerformanceTitle,

	"breaking": FieldBreakingTitle,
	"break":    FieldBreakingTitle,
	"breaks":   FieldBreakingTitle,

	"refactor":  FieldRefactorTitle,
	"refactors": FieldRefactorTitle,

	"style": FieldStyleTitle,
	"ui":    FieldStyleTitle,

	"test": FieldTestTitle,

	"docs": FieldDocsTitle,
	"doc":  FieldDocsTitle,

	"chore": FieldChoreTitle,

	"build": FieldBuildTitle,
	"ci":    FieldBuildTitle,

	"revert":   FieldRevertTitle,
	"reverts":  FieldRevertTitle,
	"reverted": FieldRevertTitle,

	"dependency":   FieldDependencyTitle,
	"dependencies": FieldDependencyTitle,
	"deps":         FieldDependencyTitle,

	"security": FieldSecurityTitle,
}

var (
	patchIters = []string{FieldFixTitle, FieldPerformanceTitle, FieldSecurityTitle}
	minorIters = []string{FieldFeatureTitle, FieldBreakingTitle, FieldRefactorTitle, FieldStyleTitle, FieldChoreTitle, FieldRevertTitle, FieldDependencyTitle}
	majorIters = []string{FieldBreakingTitle}
)

type AngularAnalyzer struct{}

type Commit struct {
	Type    string
	Scope   string
	Message string
}

func (c *Commit) Scan(message string) {
	fields := strings.Split(message, ":")
	if len(fields) == 1 {
		c.Type = "other"
		c.Message = fields[0]
		return
	}

	c.Message = fields[1]

	if strings.Contains(fields[0], "(") && strings.Contains(fields[0], ")") {
		scope := strings.Split(fields[0], "(")
		c.Type = scope[0]
		c.Scope = strings.TrimRight(scope[1], ")")
	} else {
		c.Type = fields[0]
	}

	c.Type = strings.ToLower(c.Type)
}

type info struct {
	title string
	*types.ReleaseNote
}

func (a *AngularAnalyzer) Analyze(result *types.Result) error {
	result.ReleaseNotes = make(map[string][]types.ReleaseNote)
	var wg sync.WaitGroup
	wg.Add(len(result.Commits))
	ch := make(chan info)

	for _, commit := range result.Commits {
		go func(commit *object.Commit) {
			defer wg.Done()
			title := cleanTitle(commit.Message)
			if title == "" {
				return
			}

			log.Info().Msgf("Analyzing commit %s: %s", utils.HashShort(commit.Hash), title)
			var com Commit
			com.Scan(title)

			title, ok := typeMap[com.Type]
			if !ok {
				title = FieldOtherTitle
			}

			ch <- info{
				title: title,
				ReleaseNote: &types.ReleaseNote{
					Commit: commit,
					Scope:  com.Scope,
					Desc:   com.Message,
				},
			}
		}(commit)
	}

	var majorHit, minorHit, patchHit bool
	go func() {
		for i := range ch {
			if _, ok := result.ReleaseNotes[i.title]; !ok {
				result.ReleaseNotes[i.title] = make([]types.ReleaseNote, 0)
			}

			result.ReleaseNotes[i.title] = append(result.ReleaseNotes[i.title], *i.ReleaseNote)

			if !majorHit && slices.Contains(majorIters, i.title) {
				majorHit = true
			} else if !minorHit && slices.Contains(minorIters, i.title) {
				minorHit = true
			} else if !patchHit && slices.Contains(patchIters, i.title) {
				patchHit = true
			}
		}
	}()

	wg.Wait()
	close(ch)

	if majorHit {
		result.ReleaseType = types.ReleaseTypeMajor
		result.NextRelease.Major++
	} else if minorHit {
		result.ReleaseType = types.ReleaseTypeMinor
		result.NextRelease.Minor++
	} else if patchHit {
		result.ReleaseType = types.ReleaseTypePatch
		result.NextRelease.Patch++
	}

	return nil
}

func init() {
	RegisterAnalyzer("angular", &AngularAnalyzer{})
}

func cleanTitle(message string) string {
	if strings.Contains(message, CommitIgnoreTag) {
		return ""
	}

	title := strings.Split(message, "\n")[0]
	title = strings.TrimSpace(title)

	// replace chinese characters
	title = strings.ReplaceAll(title, "（", "(")
	title = strings.ReplaceAll(title, "）", ")")
	title = strings.ReplaceAll(title, "：", ":")
	title = strings.ReplaceAll(title, "，", ",")
	title = strings.ReplaceAll(title, "。", ".")
	title = strings.ReplaceAll(title, "、", ",")
	title = strings.ReplaceAll(title, "；", ";")
	title = strings.ReplaceAll(title, "！", "!")
	title = strings.ReplaceAll(title, "？", "?")
	title = strings.ReplaceAll(title, "【", "[")
	title = strings.ReplaceAll(title, "】", "]")

	return title
}
