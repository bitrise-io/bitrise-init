package steps

import bitriseModels "github.com/bitrise-io/bitrise/models"

func RestoreGradleCache() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CacheRestoreGradleID, CacheRestoreGradleVersion)
	return stepListItem(stepIDComposite, "", "")
}

func SaveGradleCache() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CacheSaveGradleID, CacheSaveGradleVersion)
	return stepListItem(stepIDComposite, "", "")
}

func RestoreCocoapodsCache() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CacheRestoreCocoapodsID, CacheRestoreCocoapodsVersion)
	return stepListItem(stepIDComposite, "", "")
}

func SaveCocoapodsCache() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CacheSaveCocoapodsID, CacheSaveCocoapodsVersion)
	return stepListItem(stepIDComposite, "", "")
}

func RestoreCarthageCache() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CacheRestoreCarthageID, CacheRestoreCarthageVersion)
	return stepListItem(stepIDComposite, "", "")
}

func SaveCarthageCache() bitriseModels.StepListItemModel {
	stepIDComposite := stepIDComposite(CacheSaveCarthageID, CacheSaveCarthageVersion)
	return stepListItem(stepIDComposite, "", "")
}
