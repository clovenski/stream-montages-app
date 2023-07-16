using MontageJobExecutor.Models;

namespace MontageJobExecutor.Contracts
{
    public interface IMontageJobRepository
    {
        Task<UpdateMontageJobResponse> UpdateMontageJobAsync(UpdateMontageJobRequest request);
    }

    public class UpdateMontageJobRequest
    {
        public MontageJob Entity { get; set; }
    }

    public class UpdateMontageJobResponse
    {
        public MontageJob UpdatedEntity { get; set; }
    }
}
